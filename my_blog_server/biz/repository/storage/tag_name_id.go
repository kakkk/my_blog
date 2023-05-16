package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var tagNameIDStorage *TagNameIDStorage

type TagNameIDStorage struct {
	cacheX *cachex.CacheX[string, int64]
}

func GetTagNameIDStorage() *TagNameIDStorage {
	return tagNameIDStorage
}

func initTagNameIDStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[int64](ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache[int64](ctx, 1, time.Minute)
	cache := cachex.NewCacheX[string, int64]("tag_name_id", false, false).
		SetGetCacheKey(tagNameIDGetKey).
		SetGetRealData(tagNameIDGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	tagNameIDStorage = &TagNameIDStorage{
		cacheX: cache,
	}
	return nil
}

func tagNameIDGetKey(slug string) string {
	return fmt.Sprintf("tag_name_id_%v", slug)
}

func tagNameIDGetRealData(ctx context.Context, name string) (int64, error) {
	id, err := mysql.SelectTagIDByName(mysql.GetDB(ctx), name)
	if err != nil {
		return parseSqlError(id, err)
	}
	return id, nil
}

func (c *TagNameIDStorage) Get(ctx context.Context, slug string) (int64, error) {
	id, err := c.cacheX.Get(ctx, slug)
	if err != nil {
		return parseCacheXError(id, err)
	}
	return id, nil
}
