package cachex

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type CacheX[K comparable, V any] struct {
	name         string
	isSetDefault bool
	isOnlyCache  bool
	getCacheKey  func(K) string
	getRealData  func(context.Context, K) (V, error)
	mGetRealData func(context.Context, []K) (map[K]V, error)
	mGetCacheKey func([]K) []string
	cacheChain   *cacheChain[V]
	isInit       bool
}

func NewCacheX[K comparable, V any](name string, isOnlyCache bool, isSetDefault bool) *CacheX[K, V] {
	chain := newCacheChain[V]()
	return &CacheX[K, V]{
		cacheChain:   chain,
		isSetDefault: isSetDefault,
		isOnlyCache:  isOnlyCache,
		name:         name,
	}
}

func (c *CacheX[K, V]) SetGetCacheKey(f func(K) string) *CacheX[K, V] {
	if c.isInit {
		panic("SetGetCacheKey Fail, CacheX Has Init!!!")
	}
	c.getCacheKey = f
	return c
}

func (c *CacheX[K, V]) SetMGetRealData(f func(context.Context, []K) (map[K]V, error)) *CacheX[K, V] {
	if c.isInit {
		panic("SetMGetRealData Fail, CacheX Has Init!!!")
	}
	c.mGetRealData = f
	return c
}

func (c *CacheX[K, V]) SetGetRealData(f func(context.Context, K) (V, error)) *CacheX[K, V] {
	if c.isInit {
		panic("SetGetRealData Fail, CacheX Has Init!!!")
	}
	c.getRealData = f
	return c
}

func (c *CacheX[K, V]) AddCache(ctx context.Context, isSetDefault bool, cache Cache[V]) *CacheX[K, V] {
	if c.isInit {
		panic("AddCache Fail, CacheX Has Init!!!")
	}
	c.cacheChain.AddCache(ctx, c.name, isSetDefault, cache)
	return c
}

func (c *CacheX[K, V]) Initialize(ctx context.Context) error {
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
	c.mGetCacheKey = func(keys []K) []string {
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

func (c *CacheX[K, V]) setZero(ctx context.Context, key K) {
	go func() {
		cacheData := &CacheData[V]{
			CreateAt:      time.Now().UnixMilli(),
			IsDefaultData: 1,
		}
		err := c.cacheChain.Set(ctx, c.getCacheKey(key), cacheData)
		if err != nil {
			logger.Warnf(ctx, "%v set cache error: %v", c.name, err)
			return
		}
	}()
	return
}

func (c *CacheX[K, V]) Set(ctx context.Context, key K, data V) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	cacheData := &CacheData[V]{
		CreateAt: time.Now().UnixMilli(),
		Data:     data,
	}
	err := c.cacheChain.Set(ctx, c.getCacheKey(key), cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v set cache error: %v", c.name, err)
		return
	}
	return
}

func (c *CacheX[K, V]) mSetDefault(ctx context.Context, result map[K]V, keys []K) {
	go func() {
		var setDefaultKeys []K
		m := make(map[K]struct{})
		for key := range result {
			m[key] = struct{}{}
		}
		for _, key := range keys {
			if _, ok := m[key]; !ok {
				setDefaultKeys = append(setDefaultKeys, key)
			}
		}
		cacheData := make(map[string]*CacheData[V], len(setDefaultKeys))
		now := time.Now().UnixMilli()
		for _, key := range setDefaultKeys {
			cacheData[c.getCacheKey(key)] = &CacheData[V]{
				CreateAt:      now,
				IsDefaultData: 1,
			}
		}
		err := c.cacheChain.MSet(ctx, cacheData)
		if err != nil {
			logger.Warnf(ctx, "%v mset cache error: %v", c.name, err)
			return
		}
		return
	}()
}

func (c *CacheX[K, V]) MSet(ctx context.Context, data map[K]V) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	cacheData := make(map[string]*CacheData[V], len(data))
	now := time.Now().UnixMilli()
	for key, val := range data {
		cacheData[c.getCacheKey(key)] = &CacheData[V]{
			CreateAt: now,
			Data:     val,
		}
	}
	err := c.cacheChain.MSet(ctx, cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v mset cache error: %v", c.name, err)
		return
	}
	return
}

func (c *CacheX[K, V]) Get(ctx context.Context, key K) (V, error) {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	var zero V
	// 查询缓存
	fromCache, err := c.cacheChain.Get(ctx, c.getCacheKey(key))
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			logger.Warnf(ctx, "%v get cache error:[%v]", c.name, err)
		}
	}
	var data V
	// 查询到数据
	if fromCache != nil {
		// 默认值，返回NotFound
		if fromCache.IsDefault() {
			return zero, ErrNotFound
		}
		return fromCache.Data, nil
	}
	// 仅缓存模式, 直接返回
	if c.isOnlyCache {
		return zero, ErrNotFound
	}
	// 查询不到，回源
	fromDB, err := c.getRealData(ctx, key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// DB 找不到, 设置默认值
			if c.isSetDefault {
				c.setZero(ctx, key)
			}
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

func (c *CacheX[K, V]) MGet(ctx context.Context, keys []K) map[K]V {
	if !c.isInit {
		panic("CacheX not init!!!")
	}
	if c.mGetRealData == nil {
		logger.Debugf(ctx, "%v mget not set", c.name)
		return nil
	}
	if len(keys) == 0 {
		return map[K]V{}
	}
	// 去重
	keys = sliceDeduplicate(keys)
	// 查询缓存
	fromCache, err := c.cacheChain.MGet(ctx, c.mGetCacheKey(keys))
	if err != nil {
		logger.Warnf(ctx, "%v mget cache error:[%v]", c.name, err)
	}
	// 结果
	result := make(map[K]V, len(keys))
	// 需要回源的key
	var needGetRealData []K
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

		result[key] = val.Data
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
	cacheData := make(map[string]*CacheData[V], len(data))
	now := time.Now().UnixMilli()
	for key, val := range data {
		result[key] = val
		cacheData[c.getCacheKey(key)] = &CacheData[V]{
			CreateAt: now,
			Data:     val,
		}
	}
	// 设置缓存
	err = c.cacheChain.MSet(ctx, cacheData)
	if err != nil {
		logger.Warnf(ctx, "%v set cache error: %v", c.name, err)
	}
	// 设置默认值
	if c.isSetDefault && len(result) != len(keys) {
		c.mSetDefault(ctx, result, keys)
	}
	return result
}

func (c *CacheX[K, V]) Delete(ctx context.Context, key K) {
	err := c.cacheChain.Delete(ctx, c.getCacheKey(key))
	if err != nil {
		logger.Warnf(ctx, "%v delete error: %v", c.name, err)
	}
	return
}

func (c *CacheX[K, V]) MDelete(ctx context.Context, keys []K) {
	err := c.cacheChain.MDelete(ctx, c.mGetCacheKey(keys))
	if err != nil {
		logger.Warnf(ctx, "%v mdelete error: %v", c.name, err)
	}
	return
}
