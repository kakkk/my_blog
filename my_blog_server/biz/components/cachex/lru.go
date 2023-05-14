package cachex

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

type LRUCache[V any] struct {
	rw    sync.RWMutex
	cache *lru.Cache
	ttl   time.Duration
}

func NewLRUCache[V any](_ context.Context, size int, ttl time.Duration) *LRUCache[V] {
	return &LRUCache[V]{
		cache: lru.New(size),
		ttl:   ttl,
	}
}

func (l *LRUCache[V]) Get(ctx context.Context, key string) (*CacheData[V], error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	now := time.Now().UnixMilli()
	val, ok := l.cache.Get(key)
	if !ok {
		logger.Debugf(ctx, "lru cache not found, key:[%v]", key)
		return nil, ErrNotFound
	}
	result := val.(*CacheData[V])
	if l.isExpired(result.CreateAt, now) {
		logger.Debugf(ctx, "lru cache expired, key:[%v]", key)
		l.cache.Remove(key)
		return nil, ErrNotFound
	}
	return val.(*CacheData[V]), nil
}

func (l *LRUCache[V]) MGet(ctx context.Context, keys []string) (map[string]*CacheData[V], error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	now := time.Now().UnixMilli()
	result := make(map[string]*CacheData[V], len(keys))
	for _, key := range keys {
		val, ok := l.cache.Get(key)
		if !ok {
			logger.Debugf(ctx, "lru cache not found, key:[%v]", key)
			continue
		}
		res := val.(*CacheData[V])
		if l.isExpired(res.CreateAt, now) {
			logger.Debugf(ctx, "lru cache expired, key:[%v]", key)
			l.cache.Remove(key)
			continue
		}
		result[key] = res
	}
	return result, nil
}

func (l *LRUCache[V]) Set(_ context.Context, key string, data *CacheData[V]) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	now := time.Now().UnixMilli()
	data.CreateAt = now
	l.cache.Add(key, data)
	return nil
}

func (l *LRUCache[V]) MSet(_ context.Context, kvs map[string]*CacheData[V]) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	now := time.Now().UnixMilli()
	for k, v := range kvs {
		v.CreateAt = now
		l.cache.Add(k, v)
	}
	return nil
}

func (l *LRUCache[V]) Delete(_ context.Context, key string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.cache.Remove(key)
	return nil
}

func (l *LRUCache[V]) MDelete(_ context.Context, keys []string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	for _, key := range keys {
		l.cache.Remove(key)
	}
	return nil
}

func (l *LRUCache[V]) isExpired(createAt, now int64) bool {
	expire := int64(l.ttl / time.Millisecond)
	if expire <= 0 || createAt <= 0 {
		return false
	}
	if createAt+expire < now {
		return true
	}
	return false
}

func (l *LRUCache[V]) Ping(_ context.Context) (string, error) {
	if l.cache == nil {
		return "", errors.New("lru not ready")
	}
	return "Pong", nil
}

func (l *LRUCache[V]) Name() string {
	return "LRUCache"
}
