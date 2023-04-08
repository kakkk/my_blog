package cachex

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/golang/groupcache/lru"
)

type LRUCache struct {
	rw    sync.RWMutex
	cache *lru.Cache
	ttl   time.Duration
}

func NewLRUCache(ctx context.Context, size int, ttl time.Duration) *LRUCache {
	return &LRUCache{
		cache: lru.New(size),
		ttl:   ttl,
	}
}

func (l *LRUCache) Get(ctx context.Context, key string) (*CacheData, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	now := time.Now().UnixMilli()
	val, ok := l.cache.Get(key)
	if !ok {
		logger.Debugf(ctx, "lru cache not found, key:[%v]", key)
		return nil, ErrNotFound
	}
	result := val.(*CacheData)
	if l.isExpired(result.CreateAt, now) {
		logger.Debugf(ctx, "lru cache expired, key:[%v]", key)
		l.cache.Remove(key)
		return nil, ErrNotFound
	}
	return val.(*CacheData), nil
}

func (l *LRUCache) MGet(ctx context.Context, keys []string) (map[string]*CacheData, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()
	now := time.Now().UnixMilli()
	result := make(map[string]*CacheData, len(keys))
	for _, key := range keys {
		val, ok := l.cache.Get(key)
		if !ok {
			logger.Debugf(ctx, "lru cache not found, key:[%v]", key)
			continue
		}
		res := val.(*CacheData)
		if l.isExpired(res.CreateAt, now) {
			logger.Debugf(ctx, "lru cache expired, key:[%v]", key)
			l.cache.Remove(key)
			continue
		}
		result[key] = res
	}
	return result, nil
}

func (l *LRUCache) Set(ctx context.Context, key string, data *CacheData) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	now := time.Now().UnixMilli()
	data.CreateAt = now
	l.cache.Add(key, data)
	return nil
}

func (l *LRUCache) MSet(ctx context.Context, kvs map[string]*CacheData) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	now := time.Now().UnixMilli()
	for k, v := range kvs {
		v.CreateAt = now
		l.cache.Add(k, v)
	}
	return nil
}

func (l *LRUCache) Delete(ctx context.Context, key string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	l.cache.Remove(key)
	return nil
}

func (l *LRUCache) MDelete(ctx context.Context, keys []string) error {
	l.rw.Lock()
	defer l.rw.Unlock()
	for _, key := range keys {
		l.cache.Remove(key)
	}
	return nil
}

func (l *LRUCache) isExpired(createAt, now int64) bool {
	expire := int64(l.ttl / time.Millisecond)
	if expire <= 0 || createAt <= 0 {
		return false
	}
	if createAt+expire < now {
		return true
	}
	return false
}

func (l *LRUCache) Ping(ctx context.Context) (string, error) {
	if l.cache == nil {
		return "", errors.New("lru not ready")
	}
	return "Pong", nil
}

func (l *LRUCache) Name() string {
	return "LRUCache"
}
