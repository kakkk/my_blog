package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/config"
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

var categoryArticleIDsCachex *CategoryArticleIDsCachex

// CategoryArticleIDsCachex key: category_id, val: []article_id
type CategoryArticleIDsCachex struct {
	cacheX *cachex.CacheX[int64, []int64]
	expire time.Duration
}

func initCategoryArticleIDsCachex() {
	cfg := config.GetCachexSettingByName("category_article_ids")
	cache, err := infraCachex.NewCacheXBuilderByConfig[int64, []int64](context.Background(), cfg).
		SetGetRealData((&CategoryArticleIDsCachex{}).GetRealDataFn()).
		Build()

	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	categoryArticleIDsCachex = &CategoryArticleIDsCachex{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return
}

func (c *CategoryArticleIDsCachex) GetRealDataFn() func(ctx context.Context, id int64) ([]int64, error) {
	return func(ctx context.Context, id int64) ([]int64, error) {
		list, err := persistence.SelectPostIDsByCategoryID(mysql.GetDB(ctx), id)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return list, nil
	}
}

func (c *CategoryArticleIDsCachex) Get(ctx context.Context, id int64) ([]int64, error) {
	order, ok := c.cacheX.Get(ctx, id, c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}
