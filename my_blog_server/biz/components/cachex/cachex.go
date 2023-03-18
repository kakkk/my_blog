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
	Name         string
	GetCacheKey  func(KEY) string
	GetRealData  func(context.Context, KEY) (T, error)
	MGetRealData func(context.Context, []KEY) (map[KEY]T, error)
	IsSetDefault bool
	cacheChain   *cacheChain
	mGetCacheKey func([]KEY) []string
	isInit       bool
}

func NewCacheX[T Serializable[T], KEY KeyType](ctx context.Context, name string, isSetDefault bool, args ...Cache) *CacheX[T, KEY] {
	chain := newCacheChain()
	chain.AddCache(ctx, name, isSetDefault, args...)
	return &CacheX[T, KEY]{
		cacheChain: chain,
	}
}

func (c *CacheX[T, KEY]) Initialize(ctx context.Context) error {
	// 检查属性设置
	if c.Name == "" {
		return fmt.Errorf("name not set")
	}
	if c.GetCacheKey == nil {
		return fmt.Errorf("get cache key function not set")
	}
	if c.GetRealData == nil {
		return fmt.Errorf("get real data function not set")
	}
	if c.MGetRealData == nil {
		return fmt.Errorf("mget real data function not set")
	}

	// 设置mGetCacheKey
	c.mGetCacheKey = func(keys []KEY) []string {
		result := make([]string, 0, len(keys))
		for _, key := range keys {
			result = append(result, c.GetCacheKey(key))
		}
		return result
	}

	// 检查Cache是否可用
	err := c.cacheChain.CheckCache(ctx, c.Name)
	if err != nil {
		return fmt.Errorf("check cache error:[%w]", err)
	}

	c.isInit = true
	return nil
}

func (c *CacheX[T, KEY]) AddCache(ctx context.Context, isSetDefault bool, cache Cache) {
	if c.isInit {
		panic("CacheX has init!!!")
	}
	c.cacheChain.AddCache(ctx, c.Name, isSetDefault, cache)
}

func (c CacheX[T, KEY]) Set(ctx context.Context, key KEY, data T) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	cacheData := &CacheData{
		CreateAt: time.Now().UnixMilli(),
		Data:     data.Serialize(),
	}
	err := c.cacheChain.Set(ctx, c.GetCacheKey(key), cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v set cache error: %v", c.Name, err)
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
		cacheData[c.GetCacheKey(key)] = &CacheData{
			CreateAt: now,
			Data:     val.Serialize(),
		}
	}
	err := c.cacheChain.MSet(ctx, cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v `mset cache error: %v", c.Name, err)
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
	fromCache, err := c.cacheChain.Get(ctx, c.GetCacheKey(key))
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			logger.Warnf(ctx, "%v get cache error:[%v]", c.Name, err)
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
			logger.Errorf(ctx, "%v deserialize error: %v", c.Name, err)
			return zero, ErrDeserializeError
		}
		// 返回
		return data, nil
	}
	// 查询不到，回源
	fromDB, err := c.GetRealData(ctx, key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return zero, ErrNotFound
		}
		logger.Errorf(ctx, "%v get real data error: %v", c.Name, err)
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
	if len(keys) == 0 {
		return map[KEY]T{}
	}
	// 查询缓存
	fromCache, err := c.cacheChain.MGet(ctx, c.mGetCacheKey(keys))
	if err != nil {
		logger.Warnf(ctx, "%v mget cache error:[%v]", c.Name, err)
	}
	// 结果
	result := make(map[KEY]T, len(keys))
	// 需要回源的key
	var needGetRealData []KEY
	// 检查数据
	for _, key := range keys {
		val, ok := fromCache[c.GetCacheKey(key)]
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
			logger.Warnf(ctx, "%v deserialize error: %v", c.Name, err)
			needGetRealData = append(needGetRealData, key)
			continue
		}

		result[key] = data
	}
	// 全部查到，直接返回
	if len(needGetRealData) == 0 {
		return result
	}
	// 回源
	data, err := c.MGetRealData(ctx, keys)
	if err != nil && !errors.Is(err, ErrNotFound) {
		logger.Warnf(ctx, "%v get real data error: %v", c.Name, err)
		return result
	}
	// 组合结果
	cacheData := make(map[string]*CacheData, len(data))
	now := time.Now().UnixMilli()
	for key, val := range data {
		result[key] = val
		cacheData[c.GetCacheKey(key)] = &CacheData{
			CreateAt: now,
			Data:     val.Serialize(),
		}
	}
	// 设置缓存
	err = c.cacheChain.MSet(ctx, cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v set cache error: %v", c.Name, err)
	}
	return result
}

func (c *CacheX[T, KEY]) Delete(ctx context.Context, key KEY) {
	err := c.cacheChain.Delete(ctx, c.GetCacheKey(key))
	if err != nil {
		logger.Warnf(ctx, "%v delete error: %v", c.Name, err)
	}
	return
}

func (c *CacheX[T, KEY]) MDelete(ctx context.Context, keys []KEY) {
	err := c.cacheChain.MDelete(ctx, c.mGetCacheKey(keys))
	if err != nil {
		logger.Warnf(ctx, "%v mdelete error: %v", c.Name, err)
	}
	return
}
