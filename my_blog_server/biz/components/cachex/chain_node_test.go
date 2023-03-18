package cachex

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_chainNode_Delete(t *testing.T) {
	ctx := context.Background()
	cache1 := NewTestCache(ctx, 10)
	cache2 := NewTestCache(ctx, 10)
	_, tail := newTestCacheChain(cache1, cache2)

	key1 := "key1"
	key2 := "key2"
	val1 := &CacheData{
		Data: "val1",
	}
	val2 := &CacheData{
		Data: "val2",
	}
	kvs := map[string]*CacheData{
		key1: val1,
		key2: val2,
	}

	t.Run("normal", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		err := tail.Delete(ctx, key1)
		assert.Nil(t, err)
		nil1, _ := cache1.Get(ctx, key1)
		assert.Nil(t, nil1)
		nil2, _ := cache2.Get(ctx, key1)
		assert.Nil(t, nil2)
		got, _ := cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
		got, _ = cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
	})

	t.Run("delete cache 1", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.Set(ctx, key2, val2)
		err := tail.Delete(ctx, key1)
		assert.Nil(t, err)
		got, _ := cache1.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
		got, _ = cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
	})

	t.Run("delete cache 2", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.Set(ctx, key2, val2)
		_ = cache2.MSet(ctx, kvs)
		err := tail.Delete(ctx, key1)
		assert.Nil(t, err)
		got, _ := cache1.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
		got, _ = cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
	})

	t.Run("delete empty", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		err := tail.Delete(ctx, key1)
		assert.Nil(t, err)
		got, _ := cache1.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Nil(t, got)
	})

	t.Run("cache1 err", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		cacheErr := errors.New("cache1 err")
		cache1.SetError(cacheErr)
		err := tail.Delete(ctx, key1)
		assert.Nil(t, err)
		resetTestCacheErr(cache1, cache2)
		got, _ := cache1.Get(ctx, key1)
		assert.Equal(t, val1, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Equal(t, val2, got)
		got, _ = cache2.Get(ctx, key2)
		assert.Equal(t, val2, got)
	})

}

func Test_chainNode_Get(t *testing.T) {
	ctx := context.Background()
	cache1 := NewTestCache(ctx, 10)
	cache2 := NewTestCache(ctx, 10)
	head, _ := newTestCacheChain(cache1, cache2)

	key1 := "key1"
	key2 := "key2"
	val1 := &CacheData{
		Data: "val1",
	}
	val2 := &CacheData{
		Data: "val2",
	}
	kvs := map[string]*CacheData{
		key1: val1,
		key2: val2,
	}

	t.Run("all exist", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		got, err := head.Get(ctx, key1)
		assert.Nil(t, err)
		assert.EqualValues(t, val1, got)
		got, err = head.Get(ctx, key2)
		assert.Nil(t, err)
		assert.EqualValues(t, val2, got)
	})

	t.Run("cache1 exist", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		got, err := head.Get(ctx, key1)
		assert.Nil(t, err)
		assert.EqualValues(t, val1, got)
		got, err = head.Get(ctx, key2)
		assert.Nil(t, err)
		assert.EqualValues(t, val2, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key2)
		assert.Nil(t, got)
	})

	t.Run("cache2 exist", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache2.MSet(ctx, kvs)
		got, err := head.Get(ctx, key1)
		assert.Nil(t, err)
		assert.EqualValues(t, val1, got)
		got, err = head.Get(ctx, key2)
		assert.Nil(t, err)
		assert.EqualValues(t, val2, got)
		got, _ = cache1.Get(ctx, key1)
		assert.EqualValues(t, val1, got)
		got, _ = cache1.Get(ctx, key2)
		assert.EqualValues(t, val2, got)
	})

	t.Run("cache1 error & cache2 exist", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		cacheErr := errors.New("test error")
		cache1.SetError(cacheErr)
		_ = cache2.MSet(ctx, kvs)
		got, err := head.Get(ctx, key1)
		assert.Nil(t, err)
		assert.EqualValues(t, val1, got)
		got, err = head.Get(ctx, key2)
		assert.Nil(t, err)
		assert.EqualValues(t, val2, got)
		got, _ = cache1.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Nil(t, got)
	})
}

func Test_chainNode_MDelete(t *testing.T) {
	ctx := context.Background()
	cache1 := NewTestCache(ctx, 10)
	cache2 := NewTestCache(ctx, 10)
	_, tail := newTestCacheChain(cache1, cache2)

	key1 := "key1"
	key2 := "key2"
	key3 := "key3"
	val1 := &CacheData{
		Data: "val1",
	}
	val2 := &CacheData{
		Data: "val2",
	}
	val3 := &CacheData{
		Data: "val3",
	}
	kvs := map[string]*CacheData{
		key1: val1,
		key2: val2,
		key3: val3,
	}

	t.Run("normal", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		err := tail.MDelete(ctx, []string{key1, key3})
		assert.Nil(t, err)
		got, _ := cache1.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key3)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key3)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Equal(t, val2, got)
		got, _ = cache2.Get(ctx, key2)
		assert.Equal(t, val2, got)
	})

	t.Run("cache1 error", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		cache1.SetError(errors.New("test error"))
		err := tail.MDelete(ctx, []string{key1, key3})
		assert.Nil(t, err)
		cache1.ResetError()
		got, _ := cache1.Get(ctx, key1)
		assert.Equal(t, val1, got)
		got, _ = cache1.Get(ctx, key3)
		assert.Equal(t, val3, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key3)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Equal(t, val2, got)
		got, _ = cache2.Get(ctx, key2)
		assert.Equal(t, val2, got)
	})

	t.Run("cache2 error", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		cache2.SetError(errors.New("test"))
		err := tail.MDelete(ctx, []string{key1, key3})
		assert.Nil(t, err)
		cache2.ResetError()
		got, _ := cache1.Get(ctx, key1)
		assert.Nil(t, got)
		got, _ = cache1.Get(ctx, key3)
		assert.Nil(t, got)
		got, _ = cache2.Get(ctx, key1)
		assert.Equal(t, val1, got)
		got, _ = cache2.Get(ctx, key3)
		assert.Equal(t, val3, got)
		got, _ = cache1.Get(ctx, key2)
		assert.Equal(t, val2, got)
		got, _ = cache2.Get(ctx, key2)
		assert.Equal(t, val2, got)
	})

}

func Test_chainNode_MGet(t *testing.T) {
	ctx := context.Background()
	cache1 := NewTestCache(ctx, 10)
	cache2 := NewTestCache(ctx, 10)
	head, _ := newTestCacheChain(cache1, cache2)

	key1 := "key1"
	key2 := "key2"
	key3 := "key3"
	val1 := &CacheData{
		Data: "val1",
	}
	val2 := &CacheData{
		Data: "val2",
	}
	val3 := &CacheData{
		Data: "val3",
	}
	kvs := map[string]*CacheData{
		key1: val1,
		key2: val2,
		key3: val3,
	}

	t.Run("normal", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, kvs)
		_ = cache2.MSet(ctx, kvs)
		got, err := head.MGet(ctx, []string{key1, key2, key3})
		assert.Nil(t, err)
		assert.EqualValues(t, kvs, got)
	})

	t.Run("cache1 empty", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache2.MSet(ctx, kvs)
		got, err := head.MGet(ctx, []string{key1, key2, key3})
		assert.Nil(t, err)
		assert.EqualValues(t, kvs, got)
		checkCache(t, cache1, kvs)
	})

	t.Run("cache1 error", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, map[string]*CacheData{
			key1: {Data: "test"},
		})
		_ = cache2.MSet(ctx, kvs)
		cache1.SetError(errors.New("test"))
		got, err := head.MGet(ctx, []string{key1, key2, key3})
		assert.Nil(t, err)
		assert.EqualValues(t, kvs, got)
		cache1.SetError(nil)
		checkCache(t, cache1, map[string]*CacheData{
			key1: {Data: "test"},
		})
	})

	t.Run("get by cache1 cache2", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache1.MSet(ctx, map[string]*CacheData{
			key1: val1,
			key3: val3,
		})
		_ = cache2.MSet(ctx, map[string]*CacheData{
			key2: val2,
		})
		got, err := head.MGet(ctx, []string{key1, key2, key3})
		assert.Nil(t, err)
		assert.EqualValues(t, kvs, got)
		checkCache(t, cache1, kvs)
		checkCache(t, cache2, map[string]*CacheData{
			key2: val2,
		})
	})

	t.Run("all not found", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		got, err := head.MGet(ctx, []string{key1, key2, key3})
		assert.Nil(t, err)
		assert.EqualValues(t, map[string]*CacheData{}, got)
	})
	t.Run("some not found", func(t *testing.T) {
		defer resetTestCache(cache1, cache2)
		_ = cache2.MSet(ctx, map[string]*CacheData{
			key2: val2,
		})
		got, err := head.MGet(ctx, []string{key1, key2, key3})
		assert.Nil(t, err)
		assert.EqualValues(t, map[string]*CacheData{
			key2: val2,
		}, got)
		checkCache(t, cache1, map[string]*CacheData{
			key2: val2,
		})
	})
}

func Test_chainNode_MSet(t *testing.T) {

}

func Test_chainNode_Set(t *testing.T) {

}

func Test_newDefaultCacheData(t *testing.T) {
	got := newDefaultCacheData()
	assert.NotEqual(t, 0, got.CreateAt)
	assert.Equal(t, "", got.Data)
}

func newTestCacheChain(c1 Cache, c2 Cache) (*chainNode, *chainNode) {
	head := &chainNode{
		cache: &defaultCache{},
	}
	cache1 := &chainNode{
		cache: c1,
		prev:  head,
	}
	cache2 := &chainNode{
		cache: c2,
		prev:  cache1,
	}
	tail := &chainNode{
		cache: &defaultCache{},
		prev:  cache2,
	}
	head.next = cache1
	cache1.next = cache2
	cache2.next = tail
	return head, tail
}

func resetTestCache(c1 *testCache, c2 *testCache) {
	c1.ResetCache()
	c2.ResetCache()
}

func resetTestCacheErr(c1 *testCache, c2 *testCache) {
	c1.ResetError()
	c2.ResetError()
}

func checkCache(t *testing.T, c *testCache, want map[string]*CacheData) {
	ctx := context.Background()
	for k, v := range want {
		got, _ := c.Get(ctx, k)
		if v == nil {
			assert.Nil(t, got)
			continue
		}
		assert.EqualValues(t, v, got)
	}
}
