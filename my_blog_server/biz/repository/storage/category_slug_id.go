package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var categorySlugIDStorage *CategorySlugIDStorage

type CategorySlugIDStorage struct {
	cacheX *cachex.CacheX[int64, string]
}

func GetCategorySlugIDStorage() *CategorySlugIDStorage {
	return categorySlugIDStorage
}

func initCategorySlugIDStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1, time.Minute)
	cache := cachex.NewBaseValueCacheX[int64, string]("category_slug_id", false, false).
		SetGetCacheKey(categorySlugIDGetKey).
		SetGetRealData(categorySlugIDGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categorySlugIDStorage = &CategorySlugIDStorage{
		cacheX: cache,
	}
	return nil
}

func categorySlugIDGetKey(slug string) string {
	return fmt.Sprintf("category_slug_id_%v", slug)
}

func categorySlugIDGetRealData(ctx context.Context, slug string) (int64, error) {
	id, err := mysql.SelectCategoryIDBySlug(mysql.GetDB(ctx), slug)
	if err != nil {
		return parseSqlError(id, err)
	}
	return id, nil
}

func (c *CategorySlugIDStorage) Get(ctx context.Context, slug string) (int64, error) {
	id, err := c.cacheX.Get(ctx, slug)
	if err != nil {
		return parseCacheXError(id, err)
	}
	return id, nil
}
