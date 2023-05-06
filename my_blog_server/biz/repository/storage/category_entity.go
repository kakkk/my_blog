package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var categoryEntityStorage *CategoryEntityStorage

type CategoryEntityStorage struct {
	cacheX *cachex.CacheX[*entity.Category, int64]
}

func GetCategoryEntityStorage() *CategoryEntityStorage {
	return categoryEntityStorage
}

func initCategoryEntityStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1000, time.Minute)
	cache := cachex.NewCacheX[*entity.Category, int64]("category_entity", false, true).
		SetGetCacheKey(categoryStorageGetKey).
		SetGetRealData(categoryEntityStorageGetRealData).
		SetMGetRealData(categoryEntityStorageMGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categoryEntityStorage = &CategoryEntityStorage{
		cacheX: cache,
	}
	return nil
}

func categoryEntityStorageGetRealData(ctx context.Context, id int64) (*entity.Category, error) {
	db := mysql.GetDB(ctx)
	// 获取category
	category, err := mysql.SelectCategoryByID(db, id)
	if err != nil {
		if err == consts.ErrRecordNotFound {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return category, nil
}

func categoryEntityStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*entity.Category, error) {
	db := mysql.GetDB(ctx)
	category, err := mysql.MSelectCategoryByIDs(db, ids)
	if err != nil {
		if err == consts.ErrRecordNotFound {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return category, nil
}

func categoryStorageGetKey(id int64) string {
	return fmt.Sprintf("my_blog_category_entity_%v", id)
}

func (a *CategoryEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*entity.Category {
	return a.cacheX.MGet(ctx, ids)
}

func (a *CategoryEntityStorage) Get(ctx context.Context, id int64) (*entity.Category, error) {
	article, err := a.cacheX.Get(ctx, id)
	if err != nil {
		if errors.Is(err, cachex.ErrNotFound) {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("cacheX error:[%w]", err)
	}

	return article, nil
}
