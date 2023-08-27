package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
)

var categoryEntityStorage *CategoryEntityStorage

type CategoryEntityStorage struct {
	cacheX *cachex.CacheX[int64, *entity.Category]
	expire time.Duration
}

func GetCategoryEntityStorage() *CategoryEntityStorage {
	return categoryEntityStorage
}

func initCategoryEntityStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("category_entity")
	cache, err := NewCacheXBuilderByConfig[int64, *entity.Category](ctx, cfg).
		SetGetRealData(categoryEntityStorageGetRealData).
		SetMGetRealData(categoryEntityStorageMGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categoryEntityStorage = &CategoryEntityStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func categoryEntityStorageGetRealData(ctx context.Context, id int64) (*entity.Category, error) {
	db := mysql.GetDB(ctx)
	// 获取category
	category, err := mysql.SelectCategoryByID(db, id)
	if err != nil {
		return parseSqlError(category, err)
	}
	return category, nil
}

func categoryEntityStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*entity.Category, error) {
	db := mysql.GetDB(ctx)
	category, err := mysql.MSelectCategoryByIDs(db, ids)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return category, nil
}

func (a *CategoryEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*entity.Category {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *CategoryEntityStorage) Get(ctx context.Context, id int64) (*entity.Category, error) {
	article, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
