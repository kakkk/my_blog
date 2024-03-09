package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

var categorySlugIDStorage *CategorySlugIDStorage

type CategorySlugIDStorage struct {
	cacheX *cachex.CacheX[string, int64]
	expire time.Duration
}

func GetCategorySlugIDStorage() *CategorySlugIDStorage {
	return categorySlugIDStorage
}

func initCategorySlugIDStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("category_slug_id")
	cache, err := cachex2.NewCacheXBuilderByConfig[string, int64](ctx, cfg).
		SetGetRealData(categorySlugIDGetRealData).
		Build()

	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categorySlugIDStorage = &CategorySlugIDStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func categorySlugIDGetRealData(ctx context.Context, slug string) (int64, error) {
	id, err := mysql.SelectCategoryIDBySlug(mysql2.GetDB(ctx), slug)
	if err != nil {
		return parseSqlError(id, err)
	}
	return id, nil
}

func (c *CategorySlugIDStorage) Get(ctx context.Context, slug string) (int64, error) {
	id, ok := c.cacheX.Get(ctx, slug, c.expire)
	if !ok {
		return 0, consts.ErrRecordNotFound
	}
	return id, nil
}
