package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

type CategoryListCachex struct {
	cacheX *cachex.CacheX[string, []*dto.Category]
	expire time.Duration
}

var categoryListCachex *CategoryListCachex

func initCategoryListCachex() {
	cfg := config.GetCachexSettingByName("category_list")
	cx, err := cachex2.NewCacheXBuilderByConfig[string, []*dto.Category](context.Background(), cfg).
		SetGetRealData((&CategoryListCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	categoryListCachex = &CategoryListCachex{
		expire: cfg.GetExpire(),
		cacheX: cx,
	}
}

func (c *CategoryListCachex) Get(ctx context.Context) ([]*dto.Category, error) {
	got, ok := c.cacheX.Get(ctx, "", c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return got, nil
}

func (c *CategoryListCachex) GetRealDataFn() func(ctx context.Context, _ string) ([]*dto.Category, error) {
	return func(ctx context.Context, _ string) ([]*dto.Category, error) {
		db := mysql.GetDB(ctx)
		order, err := persistence.SelectCategoryOrder(db)
		if err != nil {
			return nil, fmt.Errorf("select category order error:[%w]", err)
		}
		categories, err := persistence.MSelectCategoryByIDs(db, order)
		if err != nil {
			return nil, fmt.Errorf("select category error:[%w]", err)
		}

		// 已发布的数量
		counts, err := persistence.MSelectCategoryArticleCountByCategoryIDs(db, order, true)
		if err != nil {
			return nil, fmt.Errorf("select article count error:[%w]", err)
		}

		var list []*dto.Category
		for _, id := range order {
			category, ok := categories[id]
			if !ok {
				continue
			}
			count := counts[id]
			list = append(list, &dto.Category{
				ID:           category.ID,
				CategoryName: category.CategoryName,
				Slug:         category.Slug,
				Count:        count,
			})
		}
		return list, nil
	}
}

func (c *CategoryListCachex) RebuildCache(ctx context.Context) {
	_ = c.cacheX.Delete(ctx, "")
	_, _ = c.cacheX.Get(ctx, "", 0)
}
