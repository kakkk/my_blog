package cachex

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/groupcache/lru"
)

type TestCache[T any] struct {
	rw    sync.RWMutex
	cache *lru.Cache
	err   error
}

func NewTestCache[T any](ctx context.Context, size int) *TestCache[T] {
	return &TestCache[T]{
		cache: lru.New(size),
	}
}

func (l *TestCache[T]) Get(ctx context.Context, key string) (*CacheData[T], error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	if l.err != nil {
		return nil, l.err
	}
	val, ok := l.cache.Get(key)
	if !ok {
		logger.Info(ctx, fmt.Sprintf("lru cache not found, key:[%v]", key))
		return nil, ErrNotFound
	}
	return val.(*CacheData[T]), l.err
}

func (l *TestCache[T]) MGet(ctx context.Context, keys []string) (map[string]*CacheData[T], error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	result := make(map[string]*CacheData[T], len(keys))
	if l.err != nil {
		return result, l.err
	}
	for _, key := range keys {
		val, ok := l.cache.Get(key)
		if !ok {
			logger.Info(ctx, fmt.Sprintf("lru cache not found, key:[%v]", key))
			continue
		}
		res := val.(*CacheData[T])
		result[key] = res
	}
	return result, l.err
}

func (l *TestCache[T]) Set(ctx context.Context, key string, data *CacheData[T]) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	if l.err != nil {
		return l.err
	}
	l.cache.Add(key, data)
	return l.err
}

func (l *TestCache[T]) MSet(ctx context.Context, kvs map[string]*CacheData[T]) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	if l.err != nil {
		return l.err
	}
	for k, v := range kvs {
		l.cache.Add(k, v)
	}
	return l.err
}

func (l *TestCache[T]) Delete(ctx context.Context, key string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	if l.err != nil {
		return l.err
	}
	l.cache.Remove(key)
	return l.err
}

func (l *TestCache[T]) MDelete(ctx context.Context, keys []string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	if l.err != nil {
		return l.err
	}
	for _, key := range keys {
		l.cache.Remove(key)
	}
	return l.err
}

func (l *TestCache[T]) Ping(ctx context.Context) (string, error) {
	return "Pong", nil
}

func (l *TestCache[T]) Name() string {
	return "testCache"
}

func (l *TestCache[T]) SetError(err error) {
	l.err = err
}

func (l *TestCache[T]) ResetCache() {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.cache.Clear()
	l.err = nil
}

func (l *TestCache[T]) ResetError() {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.err = nil
}
