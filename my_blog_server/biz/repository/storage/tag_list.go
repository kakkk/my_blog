package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/dto"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

type TagListStorage struct {
	cacheX *cachex.CacheX[int, *dto.TagList]
}

var tagListStorage *TagListStorage

func initTagListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[*dto.TagList](ctx, redis.GetRedisClient(ctx), 6*time.Hour)
	lruCache := cachex.NewLRUCache[*dto.TagList](ctx, 1, time.Hour)
	cache := cachex.NewCacheX[int, *dto.TagList]("tag_list", false, false).
		SetGetCacheKey(tagListGetKey).
		SetGetRealData(tagListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	tagListStorage = &TagListStorage{
		cacheX: cache,
	}
	return nil
}

func GetTagListStorage() *TagListStorage {
	return tagListStorage
}

func tagListGetKey(_ int) string {
	return "tag_list"
}

func tagListGetRealData(ctx context.Context, _ int) (*dto.TagList, error) {
	db := mysql.GetDB(ctx)
	tagEntityList, err := mysql.GetAllTag(db)
	if err != nil {
		return parseSqlError(&dto.TagList{}, err)
	}
	ids := make([]int64, 0, len(tagEntityList))
	for _, tag := range tagEntityList {
		ids = append(ids, tag.ID)
	}
	tagPostCountMap, err := mysql.MGetTagArticleCountByTagIDs(db, ids, true)
	if err != nil {
		return parseSqlError(&dto.TagList{}, err)
	}
	return dto.NewTagList(tagEntityList, tagPostCountMap), nil
}

func (c *TagListStorage) Get(ctx context.Context) (*dto.TagList, error) {
	got, err := c.cacheX.Get(ctx, 0)
	if err != nil {
		return parseCacheXError(got, err)
	}
	return got, nil
}

func (c *TagListStorage) RebuildCache(ctx context.Context) {
	c.cacheX.Delete(ctx, 0)
	_, _ = c.cacheX.Get(ctx, 0)
}
