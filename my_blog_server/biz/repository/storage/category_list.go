package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/dto"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/pkg/log"
	cachex2 "my_blog/biz/infra/repository/cachex"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

type CategoryListStorage struct {
	cacheX *cachex.CacheX[string, *dto.CategoryList]
	expire time.Duration
}

var categoryListStorage *CategoryListStorage

func initCategoryListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("category_list")
	cache, err := cachex2.NewCacheXBuilderByConfig[string, *dto.CategoryList](ctx, cfg).
		SetGetRealData(categoryListGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categoryListStorage = &CategoryListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func categoryListGetFromDB(ctx context.Context, withPublish bool) (*dto.CategoryList, error) {
	db := mysql2.GetDB(ctx)
	order, err := mysql.SelectCategoryOrder(db)
	if err != nil {
		return nil, fmt.Errorf("select category order error:[%v]", err)
	}
	categories, err := mysql.MSelectCategoryByIDs(db, order)
	if err != nil {
		return nil, fmt.Errorf("select category error:[%v]", err)
	}

	counts, err := mysql.MSelectCategoryArticleCountByCategoryIDs(db, order, withPublish)
	if err != nil {
		return nil, fmt.Errorf("select article count error:[%v]", err)
	}

	list := dto.CategoryList{}
	for _, id := range order {
		category, ok := categories[id]
		if !ok {
			log.GetLoggerWithCtx(ctx).Warnf("category not exist, category_id:[%v]", id)
			continue
		}
		count := counts[id]
		list = append(list, &dto.CategoryListItem{
			ID:    category.ID,
			Name:  category.CategoryName,
			Slug:  category.Slug,
			Count: count,
		})
	}
	return &list, nil
}

func categoryListGetRealData(ctx context.Context, _ string) (*dto.CategoryList, error) {
	return categoryListGetFromDB(ctx, true)
}

func GetCategoryListStorage() *CategoryListStorage {
	return categoryListStorage
}

func (c *CategoryListStorage) Get(ctx context.Context) (*dto.CategoryList, error) {
	got, ok := c.cacheX.Get(ctx, "", c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return got, nil
}

func (c *CategoryListStorage) GetFromDB(ctx context.Context) (*dto.CategoryList, error) {
	return categoryListGetFromDB(ctx, false)
}

func (c *CategoryListStorage) RebuildCache(ctx context.Context) {
	_ = c.cacheX.Delete(ctx, "")
	_, _ = c.cacheX.Get(ctx, "", 0)
}
