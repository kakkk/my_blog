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

type ArticleCategoriesCachex struct {
	cacheX *cachex.CacheX[int64, []*dto.Category]
	expire time.Duration
}

var articleCategories *ArticleCategoriesCachex

func initArticleCategoriesCachex() {
	cfg := config.GetCachexSettingByName("article_categories")
	cx, err := infraCachex.NewCacheXBuilderByConfig[int64, []*dto.Category](context.Background(), cfg).
		SetGetRealData((&ArticleCategoriesCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %v", err))
	}
	articleCategories = &ArticleCategoriesCachex{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
}

func (a *ArticleCategoriesCachex) GetRealDataFn() func(ctx context.Context, id int64) ([]*dto.Category, error) {
	return func(ctx context.Context, id int64) ([]*dto.Category, error) {
		db := mysql.GetDB(ctx)
		ids, err := persistence.SelectCategoryIDsByArticleID(db, id)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		categories, err := persistence.MSelectCategoryByIDs(db, ids)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return dto.NewCategoriesByModelMap(categories), nil
	}
}

func (a *ArticleCategoriesCachex) Get(ctx context.Context, id int64) ([]*dto.Category, error) {
	categories, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return categories, nil
}

func (a *ArticleCategoriesCachex) Rebuild(ctx context.Context, id int64) {
	_ = a.cacheX.Delete(ctx, id)
	_, _ = a.cacheX.Get(ctx, id, a.expire)
}
