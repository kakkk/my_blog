package cachex

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type KeyType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~string
}

type CacheX[T Serializable[T], KEY KeyType] struct {
	name         string
	isSetDefault bool
	isOnlyCache  bool
	getCacheKey  func(KEY) string
	getRealData  func(context.Context, KEY) (T, error)
	mGetRealData func(context.Context, []KEY) (map[KEY]T, error)
	mGetCacheKey func([]KEY) []string
	cacheChain   *cacheChain
	isInit       bool
}

func NewCacheX[T Serializable[T], KEY KeyType](name string, isOnlyCache bool, isSetDefault bool) *CacheX[T, KEY] {
	chain := newCacheChain()
	return &CacheX[T, KEY]{
		cacheChain:   chain,
		isSetDefault: isSetDefault,
		isOnlyCache:  isOnlyCache,
		name:         name,
	}
}

func (c *CacheX[T, KEY]) SetGetCacheKey(f func(KEY) string) *CacheX[T, KEY] {
	if c.isInit {
		panic("SetGetCacheKey Fail, CacheX Has Init!!!")
	}
	c.getCacheKey = f
	return c
}

func (c *CacheX[T, KEY]) SetMGetRealData(f func(context.Context, []KEY) (map[KEY]T, error)) *CacheX[T, KEY] {
	if c.isInit {
		panic("SetMGetRealData Fail, CacheX Has Init!!!")
	}
	c.mGetRealData = f
	return c
}

func (c *CacheX[T, KEY]) SetGetRealData(f func(context.Context, KEY) (T, error)) *CacheX[T, KEY] {
	if c.isInit {
		panic("SetGetRealData Fail, CacheX Has Init!!!")
	}
	c.getRealData = f
	return c
}

func (c *CacheX[T, KEY]) AddCache(ctx context.Context, isSetDefault bool, cache Cache) *CacheX[T, KEY] {
	if c.isInit {
		panic("AddCache Fail, CacheX Has Init!!!")
	}
	c.cacheChain.AddCache(ctx, c.name, isSetDefault, cache)
	return c
}

func (c *CacheX[T, KEY]) Initialize(ctx context.Context) error {
	// 检查属性设置
	if c.name == "" {
		return fmt.Errorf("name not set")
	}
	if c.getCacheKey == nil {
		return fmt.Errorf("get cache key function not set")
	}
	if c.getRealData == nil && !c.isOnlyCache {
		return fmt.Errorf("get real data function not set")
	}
	if c.mGetRealData == nil && !c.isOnlyCache {
		logger.Debugf(ctx, "%v mget not set", c.name)
	}

	// 设置mGetCacheKey
	c.mGetCacheKey = func(keys []KEY) []string {
		result := make([]string, 0, len(keys))
		for _, key := range keys {
			result = append(result, c.getCacheKey(key))
		}
		return result
	}

	// 检查Cache是否可用
	err := c.cacheChain.CheckCache(ctx, c.name)
	if err != nil {
		return fmt.Errorf("check cache error:[%w]", err)
	}

	c.isInit = true
	return nil
}

func (c CacheX[T, KEY]) Set(ctx context.Context, key KEY, data T) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	cacheData := &CacheData{
		CreateAt: time.Now().UnixMilli(),
		Data:     data.Serialize(),
	}
	err := c.cacheChain.Set(ctx, c.getCacheKey(key), cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v set cache error: %v", c.name, err)
		return
	}
	return
}

func (c *CacheX[T, KEY]) MSet(ctx context.Context, data map[KEY]T) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	cacheData := make(map[string]*CacheData, len(data))
	now := time.Now().UnixMilli()
	for key, val := range data {
		cacheData[c.getCacheKey(key)] = &CacheData{
			CreateAt: now,
			Data:     val.Serialize(),
		}
	}
	err := c.cacheChain.MSet(ctx, cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v `mset cache error: %v", c.name, err)
		return
	}
	return
}

func (c *CacheX[T, KEY]) Get(ctx context.Context, key KEY) (T, error) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	var zero T
	// 查询缓存
	fromCache, err := c.cacheChain.Get(ctx, c.getCacheKey(key))
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			logger.Warnf(ctx, "%v get cache error:[%v]", c.name, err)
		}
	}
	var data T
	// 查询到数据
	if fromCache != nil {
		// 默认值，返回NotFound
		if fromCache.IsDefault() {
			return zero, ErrNotFound
		}
		// 反序列化
		data, err = data.Deserialize(fromCache.Data)
		if err != nil {
			logger.Errorf(ctx, "%v deserialize error: %v", c.name, err)
			return zero, ErrDeserializeError
		}
		// 返回
		return data, nil
	}
	// 仅缓存模式, 直接返回
	if c.isOnlyCache {
		return zero, ErrNotFound
	}
	// 查询不到，回源
	fromDB, err := c.getRealData(ctx, key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return zero, ErrNotFound
		}
		logger.Errorf(ctx, "%v get real data error: %v", c.name, err)
		return zero, ErrGetRealDataError
	}
	data = fromDB
	// 设置缓存
	c.Set(ctx, key, data)
	// 返回
	return data, nil
}

func (c *CacheX[T, KEY]) MGet(ctx context.Context, keys []KEY) map[KEY]T {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	if c.mGetRealData == nil {
		logger.Debugf(ctx, "%v mget not set", c.name)
		return nil
	}
	if len(keys) == 0 {
		return map[KEY]T{}
	}
	// 查询缓存
	fromCache, err := c.cacheChain.MGet(ctx, c.mGetCacheKey(keys))
	if err != nil {
		logger.Warnf(ctx, "%v mget cache error:[%v]", c.name, err)
	}
	// 结果
	result := make(map[KEY]T, len(keys))
	// 需要回源的key
	var needGetRealData []KEY
	// 检查数据
	for _, key := range keys {
		val, ok := fromCache[c.getCacheKey(key)]
		if !ok {
			needGetRealData = append(needGetRealData, key)
			continue
		}
		// 空值
		if val.IsDefault() {
			continue
		}
		// 反序列化
		var data T
		data, err := data.Deserialize(val.Data)
		if err != nil {
			logger.Warnf(ctx, "%v deserialize error: %v", c.name, err)
			needGetRealData = append(needGetRealData, key)
			continue
		}

		result[key] = data
	}
	// 全部查到或仅缓存模式，直接返回
	if len(needGetRealData) == 0 || c.isOnlyCache {
		return result
	}
	// 回源
	data, err := c.mGetRealData(ctx, keys)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logger.Warnf(ctx, "%v get real data error: %v", c.name, err)
		return result
	}
	// 组合结果
	cacheData := make(map[string]*CacheData, len(data))
	now := time.Now().UnixMilli()
	for key, val := range data {
		result[key] = val
		cacheData[c.getCacheKey(key)] = &CacheData{
			CreateAt: now,
			Data:     val.Serialize(),
		}
	}
	// 设置缓存
	err = c.cacheChain.MSet(ctx, cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v set cache error: %v", c.name, err)
	}
	return result
}

func (c *CacheX[T, KEY]) Delete(ctx context.Context, key KEY) {
	err := c.cacheChain.Delete(ctx, c.getCacheKey(key))
	if err != nil {
		logger.Warnf(ctx, "%v delete error: %v", c.name, err)
	}
	return
}

func (c *CacheX[T, KEY]) MDelete(ctx context.Context, keys []KEY) {
	err := c.cacheChain.MDelete(ctx, c.mGetCacheKey(keys))
	if err != nil {
		logger.Warnf(ctx, "%v mdelete error: %v", c.name, err)
	}
	return
}
