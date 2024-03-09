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

var categoryPostListStorage *CategoryPostListStorage

type CategoryPostListStorage struct {
	cacheX *cachex.CacheX[int64, []int64]
	expire time.Duration
}

func GetCategoryPostListStorage() *CategoryPostListStorage {
	return categoryPostListStorage
}

func initCategoryPostListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("category_post_list")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, []int64](ctx, cfg).
		SetGetRealData(categoryPostListGetRealData).
		Build()

	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categoryPostListStorage = &CategoryPostListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func categoryPostListGetRealData(ctx context.Context, id int64) ([]int64, error) {
	list, err := mysql.SelectPostIDsByCategoryID(mysql2.GetDB(ctx), id)
	if err != nil {
		return parseSqlError(list, err)
	}
	return list, nil
}

func (p *CategoryPostListStorage) Get(ctx context.Context, id int64) ([]int64, error) {
	order, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}

// 重建缓存
func (p *CategoryPostListStorage) Rebuild(ctx context.Context, ids []int64) {
	for _, id := range ids {
		_ = p.cacheX.Delete(ctx, id)
		_, _ = p.cacheX.Get(ctx, id, 0)
	}
	return
}
