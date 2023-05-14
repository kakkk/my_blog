package cachex

import (
	"context"
	"errors"
	"time"
)

type chainNode[T any] struct {
	cache        Cache[T]
	prev         *chainNode[T]
	next         *chainNode[T]
	isSetDefault bool
}

func (c *chainNode[T]) Get(ctx context.Context, key string) (*CacheData[T], error) {
	data, err := c.cache.Get(ctx, key)
	// 其他错误打日志
	if err != nil && !errors.Is(err, ErrNotFound) {
		logger.Errorf(ctx, "%v get cache error: %v", c.cache.Name(), err)
	}

	// 有数据，直接返回
	if data != nil {
		logger.Debugf(ctx, "%v: [hit]", c.cache.Name())
		return data, nil
	}

	// 最后一级
	if c.next == nil {
		return nil, ErrNotFound
	}

	logger.Debugf(ctx, "%v: [not hit]", c.cache.Name())

	// 当前缓存查不到, 查询下一级缓存
	nextData, _ := c.next.Get(ctx, key)

	// 下一级缓存查到，更新当前缓存，返回
	if nextData != nil {
		_ = c.cache.Set(ctx, key, nextData)
		return nextData, nil
	}
	// 查不到，本级缓存插入默认值
	if c.isSetDefault {
		_ = c.cache.Set(ctx, key, newDefaultCacheData[T]())
	}
	// 返回不存在
	return nil, ErrNotFound
}

func (c *chainNode[T]) MGet(ctx context.Context, keys []string) (map[string]*CacheData[T], error) {
	data, err := c.cache.MGet(ctx, keys)
	// 其他错误打日志
	if err != nil && !errors.Is(err, ErrNotFound) {
		logger.Errorf(ctx, "%v mget cache error: %v", c.cache.Name(), err)
	}

	// 查到所有数据，直接返回
	if len(data) == len(keys) {
		return data, nil
	}

	// 最后一级
	if c.next == nil {
		return map[string]*CacheData[T]{}, ErrNotFound
	}

	// 需要从下一级缓存查询的数据
	var needGetByNext []string

	// 当前缓存不存在的数据
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			needGetByNext = append(needGetByNext, key)
		}
	}

	// 查询下一级缓存
	nextData, _ := c.next.MGet(ctx, needGetByNext)

	// 设置本级缓存
	_ = c.cache.MSet(ctx, nextData)

	// 设置默认值
	if c.isSetDefault {
		var needSetDefault []string
		for _, key := range needGetByNext {
			if _, ok := nextData[key]; !ok {
				needSetDefault = append(needSetDefault, key)
			}
		}
		defaults := make(map[string]*CacheData[T], len(needSetDefault))
		for _, key := range needSetDefault {
			defaults[key] = newDefaultCacheData[T]()
		}
		_ = c.cache.MSet(ctx, defaults)
	}

	// 组合结果返回
	for k, v := range nextData {
		data[k] = v
	}

	return data, nil
}

func (c *chainNode[T]) Set(ctx context.Context, key string, data *CacheData[T]) error {
	// 设置本级缓存
	err := c.cache.Set(ctx, key, data)
	// 有错误打日志
	if err != nil {
		logger.Errorf(ctx, "%v set cache error: %v", c.cache.Name(), err)
	}
	// 第一级
	if c.prev == nil {
		return nil
	}
	// 设置上一级缓存
	return c.prev.Set(ctx, key, data)
}

func (c *chainNode[T]) MSet(ctx context.Context, kvs map[string]*CacheData[T]) error {
	// 设置本级缓存
	err := c.cache.MSet(ctx, kvs)
	// 有错误打日志
	if err != nil {
		logger.Errorf(ctx, "%v set cache error: %v", c.cache.Name(), err)
	}
	// 第一级
	if c.prev == nil {
		return nil
	}
	// 设置上一级缓存
	return c.prev.MSet(ctx, kvs)
}

func (c *chainNode[T]) Delete(ctx context.Context, key string) error {
	// 删除本级缓存
	err := c.cache.Delete(ctx, key)
	// 有错误打日志
	if err != nil {
		logger.Errorf(ctx, "%v set cache error: %v", c.cache.Name(), err)
	}
	// 第一级
	if c.prev == nil {
		return nil
	}
	// 设置上一级缓存
	return c.prev.Delete(ctx, key)
}

func (c *chainNode[T]) MDelete(ctx context.Context, keys []string) error {
	// 删除本级缓存
	err := c.cache.MDelete(ctx, keys)
	// 有错误打日志
	if err != nil {
		logger.Errorf(ctx, "%v set cache error: %v", c.cache.Name(), err)
	}
	// 第一级
	if c.prev == nil {
		return nil
	}
	// 设置上一级缓存
	return c.prev.MDelete(ctx, keys)
}

func newDefaultCacheData[T any]() *CacheData[T] {
	var zero T
	return &CacheData[T]{
		CreateAt: time.Now().UnixMilli(),
		Data:     zero,
	}
}
