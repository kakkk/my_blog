package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/model"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

var categoryEntityStorage *CategoryEntityStorage

type CategoryEntityStorage struct {
	cacheX *cachex.CacheX[int64, *model.Category]
	expire time.Duration
}

func GetCategoryEntityStorage() *CategoryEntityStorage {
	return categoryEntityStorage
}

func initCategoryEntityStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("category_entity")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, *model.Category](ctx, cfg).
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

func categoryEntityStorageGetRealData(ctx context.Context, id int64) (*model.Category, error) {
	db := mysql2.GetDB(ctx)
	// 获取category
	category, err := mysql.SelectCategoryByID(db, id)
	if err != nil {
		return parseSqlError(category, err)
	}
	return category, nil
}

func categoryEntityStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*model.Category, error) {
	db := mysql2.GetDB(ctx)
	category, err := mysql.MSelectCategoryByIDs(db, ids)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return category, nil
}

func (a *CategoryEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*model.Category {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *CategoryEntityStorage) Get(ctx context.Context, id int64) (*model.Category, error) {
	article, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
