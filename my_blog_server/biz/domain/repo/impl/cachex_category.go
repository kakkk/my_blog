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
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

// CategoryCachex key: category_slug, val: *dto.category
type CategoryCachex struct {
	cacheX *cachex.CacheX[string, *dto.Category]
	expire time.Duration
}

var categoryCachex *CategoryCachex

func initCategoryCachex() {
	cfg := config.GetCachexSettingByName("category")
	cx, err := infraCachex.NewCacheXBuilderByConfig[string, *dto.Category](context.Background(), cfg).
		SetGetRealData((&CategoryCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	categoryCachex = &CategoryCachex{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
}

func (c *CategoryCachex) GetRealDataFn() func(ctx context.Context, slug string) (*dto.Category, error) {
	return func(ctx context.Context, slug string) (*dto.Category, error) {
		category, err := persistence.SelectCategoryBySlug(mysql.GetDB(ctx), slug)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return dto.NewCategoryByModel(category), nil
	}
}

func (c *CategoryCachex) Get(ctx context.Context, slug string) (*dto.Category, error) {
	category, ok := c.cacheX.Get(ctx, slug, c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return category, nil
}
