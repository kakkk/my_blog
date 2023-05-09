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

var categoryPostListStorage *CategoryPostListStorage

type CategoryPostListStorage struct {
	cacheX *cachex.CacheX[*dto.Int64List, int64]
}

func GetCategoryPostListStorage() *CategoryPostListStorage {
	return categoryPostListStorage
}

func initCategoryPostListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1, time.Minute)
	cache := cachex.NewSerializableCacheX[*dto.Int64List, int64]("category_post_list", false, false).
		SetGetCacheKey(categoryPostListGetKey).
		SetGetRealData(categoryPostListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categoryPostListStorage = &CategoryPostListStorage{
		cacheX: cache,
	}
	return nil
}

func categoryPostListGetKey(id int64) string {
	return fmt.Sprintf("category_post_list_%v", id)
}

func categoryPostListGetRealData(ctx context.Context, id int64) (*dto.Int64List, error) {
	list, err := mysql.SelectPostIDsByCategoryID(mysql.GetDB(ctx), id)
	if err != nil {
		return parseSqlError[*dto.Int64List](err)
	}
	return dto.NewInt64List(list), nil
}

func (p *CategoryPostListStorage) Get(ctx context.Context, id int64) ([]int64, error) {
	order, err := p.cacheX.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return order.ToInt64List(), err
}

// 重建缓存
func (p *CategoryPostListStorage) Rebuild(ctx context.Context, ids []int64) {
	for _, id := range ids {
		p.cacheX.Delete(ctx, id)
		_, _ = p.cacheX.Get(ctx, id)
	}
	return
}
