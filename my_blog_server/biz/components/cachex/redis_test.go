package cachex

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisCache(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// 成功
	t.Run("success", func(t *testing.T) {
		defer resetMiniRedis(mr)
		got := NewRedisCache[string](ctx, client, time.Second*30)
		assert.NotNil(t, got)
	})

	// 失败
	t.Run("fail", func(t *testing.T) {
		defer resetMiniRedis(mr)
		assert.Panics(t, func() {
			mr.SetError("test")
			_ = NewRedisCache[string](ctx, client, time.Second*30)
		})
	})
}

func TestRedisCache_Delete(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &RedisCache[string]{
		redisClient: client,
		ttl:         time.Second * 30,
	}

	// 成功
	t.Run("success", func(t *testing.T) {
		defer resetMiniRedis(mr)
		value := getValue("test", time.Now())
		key := "key"
		_ = mr.Set(key, value)
		err := cache.Delete(ctx, key)
		assert.Nil(t, err)
		_, err = mr.Get(key)
		assert.ErrorIs(t, err, miniredis.ErrKeyNotFound)
	})

	// 失败
	t.Run("fail", func(t *testing.T) {
		defer resetMiniRedis(mr)
		mr.SetError("test")
		err := cache.Delete(ctx, "key")
		assert.ErrorIs(t, err, ErrCacheError)
	})

}

func TestRedisCache_Get(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &RedisCache[string]{
		redisClient: client,
		ttl:         time.Second * 30,
	}

	// 成功: 获取一个
	t.Run("success/single", func(t *testing.T) {
		defer resetMiniRedis(mr)
		value := getValue("test", time.Now())
		key := "key"
		_ = mr.Set(key, value)
		got, err := cache.Get(ctx, key)
		assert.Nil(t, err)
		assert.Equal(t, "test", got.Data)
	})

	// 成功: 获取多个
	t.Run("success/multi", func(t *testing.T) {
		defer resetMiniRedis(mr)
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for k, v := range kvs {
			_ = mr.Set(k, getValue(v, time.Now()))
		}
		for k, v := range kvs {
			got, err := cache.Get(ctx, k)
			assert.Nil(t, err)
			assert.Equal(t, v, got.Data)
		}
	})

	// 失败: 不存在
	t.Run("fail/not_found", func(t *testing.T) {
		defer resetMiniRedis(mr)
		got, err := cache.Get(ctx, "key")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, got)
	})

	// 失败: 数据已过期
	t.Run("fail/data_expired", func(t *testing.T) {
		defer resetMiniRedis(mr)
		value := getValue("test", time.Now().Add(-time.Hour))
		key := "key"
		_ = mr.Set(key, value)
		got, err := cache.Get(ctx, key)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, got)
	})

	// 失败: Redis已过期
	t.Run("fail/redis_expired", func(t *testing.T) {
		defer resetMiniRedis(mr)
		value := getValue("test", time.Now())
		key := "key"
		_ = mr.Set(key, value)
		mr.SetTTL(key, cache.ttl)
		mr.FastForward(time.Minute)
		got, err := cache.Get(ctx, key)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Nil(t, got)
	})

	// 失败: json错误
	t.Run("fail/json_error", func(t *testing.T) {
		defer resetMiniRedis(mr)
		key := "key"
		_ = mr.Set(key, "{1")
		got, err := cache.Get(ctx, key)
		assert.ErrorIs(t, err, ErrCacheError)
		assert.Nil(t, got)
	})

	// 失败: redis错误
	t.Run("fail/redis_error", func(t *testing.T) {
		defer resetMiniRedis(mr)
		mr.SetError("test")
		got, err := cache.Get(ctx, "key")
		assert.ErrorIs(t, err, ErrCacheError)
		assert.Nil(t, got)
	})
}

func TestRedisCache_MDelete(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &RedisCache[string]{
		redisClient: client,
		ttl:         time.Second * 30,
	}

	// 成功: 单个
	t.Run("success/single", func(t *testing.T) {
		defer resetMiniRedis(mr)
		value := getValue("test", time.Now())
		key := "key"
		_ = mr.Set(key, value)
		err := cache.MDelete(ctx, []string{key})
		assert.Nil(t, err)
		_, err = mr.Get(key)
		assert.ErrorIs(t, err, miniredis.ErrKeyNotFound)
	})

	// 成功: 多个
	t.Run("success/multi", func(t *testing.T) {
		defer resetMiniRedis(mr)
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for k, v := range kvs {
			_ = mr.Set(k, getValue(v, time.Now()))
		}
		err := cache.MDelete(ctx, []string{"key1", "key2", "key3"})
		assert.Nil(t, err)
		for k, _ := range kvs {
			_, err = mr.Get(k)
			assert.ErrorIs(t, err, miniredis.ErrKeyNotFound)
		}
	})

	// 失败
	t.Run("fail", func(t *testing.T) {
		defer resetMiniRedis(mr)
		defer mr.SetError("")
		mr.SetError("test")
		err := cache.MDelete(ctx, []string{"key1", "key2", "key3"})
		assert.ErrorIs(t, err, ErrCacheError)
	})

}

func TestRedisCache_MGet(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &RedisCache[string]{
		redisClient: client,
		ttl:         time.Second * 30,
	}

	// 成功: 获取一个
	t.Run("success/single", func(t *testing.T) {
		defer resetMiniRedis(mr)
		value := getValue("test", time.Now())
		key := "key"
		_ = mr.Set(key, value)
		got, err := cache.Get(ctx, key)
		assert.Nil(t, err)
		assert.Equal(t, "test", got.Data)
	})

	// 成功: 获取多个
	t.Run("success/multi", func(t *testing.T) {
		defer resetMiniRedis(mr)
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for k, v := range kvs {
			_ = mr.Set(k, getValue(v, time.Now()))
		}
		got, err := cache.MGet(ctx, []string{"key1", "key2", "key3"})
		assert.Nil(t, err)
		for k, v := range got {
			assert.Equal(t, kvs[k], v.Data)
		}
	})

	// 成功: 部分不存在
	t.Run("success/multi/some_not_exist", func(t *testing.T) {
		defer resetMiniRedis(mr)
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for k, v := range kvs {
			_ = mr.Set(k, getValue(v, time.Now()))
		}
		got, err := cache.MGet(ctx, []string{"key1", "key2", "key3", "key4"})
		assert.Nil(t, err)
		assert.Equal(t, 3, len(got))
		for k, v := range got {
			assert.Equal(t, kvs[k], v.Data)
		}
	})

	// 成功: 部分过期
	t.Run("success/multi/some_expired", func(t *testing.T) {
		defer resetMiniRedis(mr)
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for k, v := range kvs {
			_ = mr.Set(k, getValue(v, time.Now()))
		}
		_ = mr.Set("key4", getValue("value4", time.Now().Add(-time.Hour)))
		got, err := cache.MGet(ctx, []string{"key1", "key2", "key3", "key4"})
		assert.Nil(t, err)
		assert.Equal(t, 3, len(got))
		for k, v := range got {
			assert.Equal(t, kvs[k], v.Data)
		}
	})

	// 成功: 所有情况
	t.Run("success/multi/all_case", func(t *testing.T) {
		defer resetMiniRedis(mr)
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for k, v := range kvs {
			_ = mr.Set(k, getValue(v, time.Now()))
		}
		_ = mr.Set("key4", getValue("value4", time.Now().Add(-time.Hour)))
		_ = mr.Set("key5", "{1")
		got, err := cache.MGet(ctx, []string{"key1", "key2", "key3", "key4", "key5", "key6"})
		assert.Nil(t, err)
		assert.Equal(t, 3, len(got))
		for k, v := range got {
			assert.Equal(t, kvs[k], v.Data)
		}
	})

	// 失败
	t.Run("fail", func(t *testing.T) {
		defer resetMiniRedis(mr)
		mr.SetError("test")
		got, err := cache.MGet(ctx, []string{"key1", "key2"})
		assert.ErrorIs(t, err, ErrCacheError)
		assert.Nil(t, got)
	})

}

func TestRedisCache_MSet(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &RedisCache[string]{
		redisClient: client,
		ttl:         time.Second * 30,
	}

	// 成功: 单个
	t.Run("success/single", func(t *testing.T) {
		defer resetMiniRedis(mr)
		now := time.Now()
		data := &CacheData[string]{
			CreateAt: now.UnixMilli(),
			Data:     "test",
		}
		key := "key"
		err := cache.MSet(ctx, map[string]*CacheData[string]{key: data})
		assert.Nil(t, err)
		val, _ := mr.Get(key)
		assert.JSONEq(t, getValue("test", now), val)
	})

	// 成功: 多个
	t.Run("success/multi", func(t *testing.T) {
		defer resetMiniRedis(mr)
		now := time.Now()
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		data := make(map[string]*CacheData[string], 3)
		for k, v := range kvs {
			data[k] = &CacheData[string]{
				CreateAt: now.UnixMilli(),
				Data:     v,
			}
		}
		err := cache.MSet(ctx, data)
		assert.Nil(t, err)
		for k, v := range kvs {
			val, _ := mr.Get(k)
			assert.JSONEq(t, getValue(v, now), val)
		}
	})

	// 失败
	t.Run("fail", func(t *testing.T) {
		defer resetMiniRedis(mr)
		now := time.Now()
		kvs := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		data := make(map[string]*CacheData[string], 3)
		for k, v := range kvs {
			data[k] = &CacheData[string]{
				CreateAt: now.UnixMilli(),
				Data:     v,
			}
		}
		mr.SetError("test")
		err := cache.MSet(ctx, data)
		assert.ErrorIs(t, err, ErrCacheError)
		for k := range kvs {
			_, err := mr.Get(k)
			assert.ErrorIs(t, err, miniredis.ErrKeyNotFound)
		}
	})
}

func TestRedisCache_Set(t *testing.T) {
	ctx := context.Background()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cache := &RedisCache[string]{
		redisClient: client,
		ttl:         time.Second * 30,
	}

	// 成功
	t.Run("success", func(t *testing.T) {
		defer resetMiniRedis(mr)
		now := time.Now()
		data := &CacheData[string]{
			CreateAt: now.UnixMilli(),
			Data:     "test",
		}
		key := "key"
		err := cache.Set(ctx, key, data)
		assert.Nil(t, err)
		val, _ := mr.Get(key)
		assert.JSONEq(t, getValue("test", now), val)
	})

	// 成功/测试过期
	t.Run("success/test_expired", func(t *testing.T) {
		defer resetMiniRedis(mr)
		now := time.Now()
		data := &CacheData[string]{
			CreateAt: now.UnixMilli(),
			Data:     "test",
		}
		key := "key"
		err := cache.Set(ctx, key, data)
		assert.Nil(t, err)
		val, _ := mr.Get(key)
		assert.JSONEq(t, getValue("test", now), val)
		mr.FastForward(time.Minute)
		_, err = mr.Get(key)
		assert.ErrorIs(t, err, miniredis.ErrKeyNotFound)
	})

	// 失败
	t.Run("fail", func(t *testing.T) {
		defer resetMiniRedis(mr)
		mr.SetError("test")
		now := time.Now()
		data := &CacheData[string]{
			CreateAt: now.UnixMilli(),
			Data:     "test",
		}
		key := "key"
		err := cache.Set(ctx, key, data)
		assert.ErrorIs(t, err, ErrCacheError)
		_, err = mr.Get(key)
		assert.ErrorIs(t, err, miniredis.ErrKeyNotFound)
	})
}

func TestRedisCache_marshal(t *testing.T) {
	cache := &RedisCache[string]{}

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		data := &CacheData[string]{
			CreateAt: now.UnixMilli(),
			Data:     "test",
		}
		got := cache.marshal(data)
		assert.JSONEq(t, getValue("test", now), got)
	})

	t.Run("data_is_nil", func(t *testing.T) {
		got := cache.marshal(nil)
		assert.Equal(t, "", got)
	})
}

func TestRedisCache_unmarshal(t *testing.T) {
	cache := &RedisCache[string]{}

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		want := &CacheData[string]{
			CreateAt: now.UnixMilli(),
			Data:     "test",
		}
		got, err := cache.unmarshal(getValue("test", now))
		assert.Nil(t, err)
		assert.EqualValues(t, want, got)
	})

	t.Run("fail", func(t *testing.T) {
		got, err := cache.unmarshal("{1")
		assert.NotNil(t, err)
		assert.Nil(t, got)
	})
}

func TestRedisCache_isExpired(t *testing.T) {
	t.Run("expired", func(t *testing.T) {
		now := time.Now()
		cache := &RedisCache[string]{
			ttl: 30 * time.Second,
		}
		got := cache.isExpired(now.Add(-31*time.Second).UnixMilli(), now.UnixMilli())
		assert.True(t, got)
	})
	t.Run("not_expired/less", func(t *testing.T) {
		now := time.Now()
		cache := &RedisCache[string]{
			ttl: 30 * time.Second,
		}
		got := cache.isExpired(now.Add(-29*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
	})
	t.Run("not_expired/equal", func(t *testing.T) {
		now := time.Now()
		cache := &RedisCache[string]{
			ttl: 30 * time.Second,
		}
		got := cache.isExpired(now.Add(-30*time.Second).UnixMilli(), now.UnixMilli())
		assert.False(t, got)
	})
	t.Run("never_expired", func(t *testing.T) {
		now := time.Now()
		cache := &RedisCache[string]{
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

func TestRedisCache_Name(t *testing.T) {
	cache := &RedisCache[string]{}
	assert.Equal(t, "RedisCache", cache.Name())
}

func getValue(data string, createAt time.Time) string {
	return fmt.Sprintf(`{"c":%v,"d":"%v","z":0}`, cast.ToString(createAt.UnixMilli()), data)
}

func resetMiniRedis(mr *miniredis.Miniredis) {
	mr.FlushAll()
	mr.SetError("")
}
