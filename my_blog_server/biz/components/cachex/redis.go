package cachex

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisCache struct {
	redisClient *redis.Client
	ttl         time.Duration
}

func NewRedisCache(ctx context.Context, client *redis.Client, ttl time.Duration) *RedisCache {
	_, err := client.Ping().Result()
	if err != nil {
		panic("redis client not ready")
	}
	return &RedisCache{
		redisClient: client,
		ttl:         ttl,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (*CacheData, error) {
	now := time.Now().UnixMilli()
	val, err := r.redisClient.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Debugf(ctx, "redis cache not found, key:[%v]", key)
			return nil, ErrNotFound
		}
		logger.Errorf(ctx, "redis cache get fail, key[%v], error:[%v]", key, err)
		return nil, ErrCacheError
	}
	data, err := r.unmarshal(val)
	if err != nil {
		logger.Errorf(ctx, "redis cache unmarshal fail, key[%v], error:[%v]", key, err)
		return nil, ErrCacheError
	}
	if r.isExpired(data.CreateAt, now) {
		logger.Debugf(ctx, "redis cache expired, key:[%v]", key)
		_ = r.Delete(ctx, key)
		return nil, ErrNotFound
	}

	return data, nil
}

func (r *RedisCache) MGet(ctx context.Context, keys []string) (map[string]*CacheData, error) {
	now := time.Now().UnixMilli()
	values, err := r.redisClient.MGet(keys...).Result()
	if err != nil {
		logger.Errorf(ctx, "redis cache mget fail, error:[%v]", err)
		return nil, ErrCacheError
	}
	result := make(map[string]*CacheData, len(keys))
	for i, key := range keys {
		val, ok := values[i].(string)
		if !ok {
			logger.Debugf(ctx, "redis cache not found, key:[%v]", key)
			continue
		}
		res, err := r.unmarshal(val)
		if err != nil {
			logger.Errorf(ctx, "redis cache unmarshal fail, key[%v], error:[%v]", key, err)
			continue
		}
		if r.isExpired(res.CreateAt, now) {
			logger.Debugf(ctx, "redis cache expired, key:[%v]", key)
			_ = r.Delete(ctx, key)
			continue
		}
		result[key] = res
	}
	return result, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, data *CacheData) error {
	err := r.redisClient.Set(key, r.marshal(data), r.ttl).Err()
	if err != nil {
		logger.Errorf(ctx, "redis cache set fail, key:[%v], error:[%v]", key, err)
		return ErrCacheError
	}
	return nil
}

func (r *RedisCache) MSet(ctx context.Context, kvs map[string]*CacheData) error {
	pipe := r.redisClient.Pipeline()
	defer func(pipe redis.Pipeliner) {
		err := pipe.Close()
		if err != nil {
			logger.Errorf(ctx, "redis pipeline close fail, error:[%v]", err)
		}
	}(pipe)
	for k, v := range kvs {
		pipe.Set(k, r.marshal(v), r.ttl)
	}
	_, err := pipe.Exec()
	if err != nil {
		logger.Errorf(ctx, "redis pipeline exec fail, error:[%v]", err)
		return ErrCacheError
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.redisClient.Del(key).Err()
	if err != nil {
		logger.Errorf(ctx, "redis delete fail, error:[%v]", err)
		return ErrCacheError
	}
	return nil
}

func (r *RedisCache) MDelete(ctx context.Context, keys []string) error {
	err := r.redisClient.Del(keys...).Err()
	if err != nil {
		logger.Errorf(ctx, "redis delete fail, error:[%v]", err)
		return ErrCacheError
	}
	return nil
}

func (r *RedisCache) unmarshal(data string) (*CacheData, error) {
	result := &CacheData{}
	err := json.Unmarshal([]byte(data), result)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error:[%v]", err)
	}
	return result, nil
}

func (r *RedisCache) marshal(data *CacheData) string {
	if data == nil {
		return ""
	}
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

func (r *RedisCache) isExpired(createAt, now int64) bool {
	expire := int64(r.ttl / time.Millisecond)
	if expire <= 0 || createAt <= 0 {
		return false
	}
	if createAt+expire < now {
		return true
	}
	return false
}

func (r *RedisCache) Ping(ctx context.Context) (string, error) {
	return r.redisClient.Ping().Result()
}

func (r *RedisCache) Name() string {
	return "RedisCache"
}
