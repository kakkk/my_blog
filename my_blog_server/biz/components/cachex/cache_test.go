package cachex

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/groupcache/lru"
)

type testCache struct {
	rw    sync.RWMutex
	cache *lru.Cache
	err   error
}

func NewTestCache(ctx context.Context, size int) *testCache {
	return &testCache{
		cache: lru.New(size),
	}
}

func (l *testCache) Get(ctx context.Context, key string) (*CacheData, error) {
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
	return val.(*CacheData), l.err
}

func (l *testCache) MGet(ctx context.Context, keys []string) (map[string]*CacheData, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	result := make(map[string]*CacheData, len(keys))
	if l.err != nil {
		return result, l.err
	}
	for _, key := range keys {
		val, ok := l.cache.Get(key)
		if !ok {
			logger.Info(ctx, fmt.Sprintf("lru cache not found, key:[%v]", key))
			continue
		}
		res := val.(*CacheData)
		result[key] = res
	}
	return result, l.err
}

func (l *testCache) Set(ctx context.Context, key string, data *CacheData) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	if l.err != nil {
		return l.err
	}
	l.cache.Add(key, data)
	return l.err
}

func (l *testCache) MSet(ctx context.Context, kvs map[string]*CacheData) error {
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

func (l *testCache) Delete(ctx context.Context, key string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	if l.err != nil {
		return l.err
	}
	l.cache.Remove(key)
	return l.err
}

func (l *testCache) MDelete(ctx context.Context, keys []string) error {
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

func (l *testCache) Ping(ctx context.Context) (string, error) {
	return "Pong", nil
}

func (l *testCache) Name() string {
	return "testCache"
}

func (l *testCache) SetError(err error) {
	l.err = err
}

func (l *testCache) ResetCache() {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.cache.Clear()
	l.err = nil
}

func (l *testCache) ResetError() {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.err = nil
}
