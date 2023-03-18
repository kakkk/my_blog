package cachex

import (
	"context"
	"testing"
	"time"

	"github.com/golang/groupcache/lru"
	"github.com/stretchr/testify/assert"
)

func TestNewLRUCache(t *testing.T) {
	ctx := context.Background()
	t.Run("normal", func(t *testing.T) {
		got := NewLRUCache(ctx, 10, time.Second*10)
		assert.Equal(t, time.Second*10, got.ttl)
		assert.NotNil(t, got)
		assert.NotNil(t, got.cache)
	})
}

func TestLRUCache_Delete(t *testing.T) {
	ctx := context.Background()
	t.Run("normal", func(t *testing.T) {
		key, val := "key", "val"
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Second,
		}
		c.Add(key, val)
		err := cache.Delete(ctx, key)
		assert.Nil(t, err)
		got, ok := c.Get(key)
		assert.Nil(t, got)
		assert.False(t, ok)
	})
}

func TestLRUCache_Get(t *testing.T) {
	ctx := context.Background()
	t.Run("normal", func(t *testing.T) {
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Hour,
		}
		key := "key"
		val := &CacheData{
			CreateAt: time.Now().UnixMilli(),
			Data:     "test",
		}
		c.Add(key, val)
		got, err := cache.Get(ctx, key)
		assert.Nil(t, err)
		assert.EqualValues(t, val, got)
	})

	t.Run("fail/not_exist", func(t *testing.T) {
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Hour,
		}
		got, err := cache.Get(ctx, "key")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, got)
	})

	t.Run("fail/expired", func(t *testing.T) {
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Second,
		}
		key := "key"
		val := &CacheData{
			CreateAt: time.Now().Add(-1 * time.Hour).UnixMilli(),
			Data:     "test",
		}
		c.Add(key, val)
		got, err := cache.Get(ctx, "key")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, got)
		_, ok := c.Get(key)
		assert.False(t, ok)
	})

}

func TestLRUCache_MDelete(t *testing.T) {
	ctx := context.Background()
	t.Run("normal", func(t *testing.T) {
		c := lru.New(5)
		now := time.Now()
		kvs := map[string]*CacheData{
			"key1": {CreateAt: now.UnixMilli(), Data: "value1"},
			"key2": {CreateAt: now.UnixMilli(), Data: "value2"},
			"key3": {CreateAt: now.UnixMilli(), Data: "value3"},
		}
		cache := &LRUCache{
			cache: c,
			ttl:   time.Second,
		}
		for k, v := range kvs {
			c.Add(k, v)
		}
		err := cache.MDelete(ctx, []string{"key1", "key3"})
		assert.Nil(t, err)
		_, ok := c.Get("key1")
		assert.False(t, ok)
		val2, ok := c.Get("key2")
		assert.True(t, ok)
		assert.Equal(t, kvs["key2"], val2)
		_, ok = c.Get("key3")
		assert.False(t, ok)
	})
}

func TestLRUCache_MGet(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		now := time.Now()
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Hour,
		}
		kvs := map[string]*CacheData{
			"key1": {CreateAt: now.UnixMilli(), Data: "value1"},
			"key2": {CreateAt: now.UnixMilli(), Data: "value2"},
			"key3": {CreateAt: now.UnixMilli(), Data: "value3"},
		}
		for k, v := range kvs {
			c.Add(k, v)
		}
		got, err := cache.MGet(ctx, []string{"key1", "key2", "key3"})
		assert.Nil(t, err)
		assert.EqualValues(t, kvs, got)
	})

	t.Run("all_cases", func(t *testing.T) {
		now := time.Now()
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Hour,
		}
		kvs := map[string]*CacheData{
			"key1": {CreateAt: now.UnixMilli(), Data: "value1"},
			"key2": {CreateAt: now.UnixMilli(), Data: "value2"},
			"key3": {CreateAt: now.UnixMilli(), Data: "value3"},
			"key4": {CreateAt: now.Add(-2 * time.Hour).UnixMilli(), Data: "value4"},
		}
		for k, v := range kvs {
			c.Add(k, v)
		}
		want := map[string]*CacheData{
			"key1": {CreateAt: now.UnixMilli(), Data: "value1"},
			"key2": {CreateAt: now.UnixMilli(), Data: "value2"},
			"key3": {CreateAt: now.UnixMilli(), Data: "value3"},
		}
		got, err := cache.MGet(ctx, []string{"key1", "key2", "key3", "key4", "key5"})
		assert.Nil(t, err)
		assert.EqualValues(t, want, got)
		length := c.Len()
		assert.Equal(t, 3, length)
	})
}

func TestLRUCache_MSet(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		now := time.Now()
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Hour,
		}
		kvs := map[string]*CacheData{
			"key1": {CreateAt: now.UnixMilli(), Data: "value1"},
			"key2": {CreateAt: now.UnixMilli(), Data: "value2"},
			"key3": {CreateAt: now.UnixMilli(), Data: "value3"},
		}
		err := cache.MSet(ctx, kvs)
		assert.Nil(t, err)
		for k, v := range kvs {
			val, ok := c.Get(k)
			assert.True(t, ok)
			assert.EqualValues(t, v, val)
		}
	})
}

func TestLRUCache_Set(t *testing.T) {
	ctx := context.Background()

	t.Run("normal", func(t *testing.T) {
		c := lru.New(5)
		cache := &LRUCache{
			cache: c,
			ttl:   time.Hour,
		}
		key := "key"
		val := &CacheData{
			CreateAt: time.Now().UnixMilli(),
			Data:     "test",
		}
		err := cache.Set(ctx, key, val)
		assert.Nil(t, err)
		got, ok := c.Get(key)
		assert.True(t, ok)
		assert.EqualValues(t, val, got)
	})
}

func TestLRUCache_isExpired(t *testing.T) {
	t.Run("expired", func(t *testing.T) {
		now := time.Now()
		cache := &LRUCache{
			ttl: 30 * time.Second,
		}
		got := cache.isExpired(now.Add(-31*time.Second).UnixMilli(), now.UnixMilli())
		assert.True(t, got)
	})
	t.Run("not_expired/less", func(t *testing.T) {
		now := time.Now()
		cache := &LRUCache{
			ttl: 30 * time.Second,
		}
		got := cache.isExpired(now.Add(-29*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
	})
	t.Run("not_expired/equal", func(t *testing.T) {
		now := time.Now()
		cache := &LRUCache{
			ttl: 30 * time.Second,
		}
		got := cache.isExpired(now.Add(-30*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
	})
	t.Run("never_expired", func(t *testing.T) {
		now := time.Now()
		cache := &LRUCache{
			ttl: 0,
		}
		got := cache.isExpired(now.Add(-30*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
		got = cache.isExpired(now.Add(-29*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
		got = cache.isExpired(now.Add(-31*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
		got = cache.isExpired(now.UnixMilli(), now.UnixMilli())
		assert.False(t, got)
	})
}
