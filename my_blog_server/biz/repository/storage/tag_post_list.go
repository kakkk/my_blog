package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var tagPostListStorage *TagPostListStorage

type TagPostListStorage struct {
	cacheX *cachex.CacheX[int64, []int64]
}

func GetTagPostListStorage() *TagPostListStorage {
	return tagPostListStorage
}

func initTagPostListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[[]int64](ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache[[]int64](ctx, 1, time.Minute)
	cache := cachex.NewCacheX[int64, []int64]("tag_post_list", false, false).
		SetGetCacheKey(tagPostListGetKey).
		SetGetRealData(tagPostListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	tagPostListStorage = &TagPostListStorage{
		cacheX: cache,
	}
	return nil
}

func tagPostListGetKey(id int64) string {
	return fmt.Sprintf("tag_post_list_%v", id)
}

func tagPostListGetRealData(ctx context.Context, id int64) ([]int64, error) {
	list, err := mysql.SelectPostIDsByTagID(mysql.GetDB(ctx), id)
	if err != nil {
		return parseSqlError(list, err)
	}
	return list, nil
}

func (p *TagPostListStorage) Get(ctx context.Context, id int64) ([]int64, error) {
	order, err := p.cacheX.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return order, err
}

// 重建缓存
func (p *TagPostListStorage) Rebuild(ctx context.Context, ids []int64) {
	for _, id := range ids {
		p.cacheX.Delete(ctx, id)
		_, _ = p.cacheX.Get(ctx, id)
	}
	return
}
