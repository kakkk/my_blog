package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var postCategoryListStorage *PostCategoryListStorage

type PostCategoryListStorage struct {
	cacheX *cachex.CacheX[int64, []*entity.Category]
}

func GetPostCategoryListStorage() *PostCategoryListStorage {
	return postCategoryListStorage
}

func initPostCategoryListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[[]*entity.Category](ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache[[]*entity.Category](ctx, 1, time.Minute)
	cache := cachex.NewCacheX[int64, []*entity.Category]("post_category_list", false, true).
		SetGetCacheKey(postCategoryListGetKey).
		SetGetRealData(postCategoryListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postCategoryListStorage = &PostCategoryListStorage{
		cacheX: cache,
	}
	return nil
}

func postCategoryListGetKey(id int64) string {
	return fmt.Sprintf("post_category_list_%v", id)
}

func postCategoryListGetRealData(ctx context.Context, id int64) ([]*entity.Category, error) {
	db := mysql.GetDB(ctx)
	categoryIDs, err := mysql.SelectCategoryIDsByArticleID(db, id)
	if err != nil {
		return parseSqlError([]*entity.Category{}, err)
	}
	cMap := GetCategoryEntityStorage().MGet(ctx, categoryIDs)
	result := make([]*entity.Category, 0, len(cMap))
	for _, category := range cMap {
		result = append(result, category)
	}
	return result, nil
}

func (p *PostCategoryListStorage) Get(ctx context.Context, id int64) ([]*entity.Category, error) {
	list, err := p.cacheX.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return list, err
}
