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

type ArticlePostIDsCachex struct {
	cacheX *cachex.CacheX[string, []int64]
	expire time.Duration
}

var articlePostIDsCachex *ArticlePostIDsCachex

func initArticlePostIDsCachex() {
	cfg := config.GetCachexSettingByName("article_post_ids")
	cache, err := infraCachex.NewCacheXBuilderByConfig[string, []int64](context.Background(), cfg).
		SetGetRealData((&ArticlePostIDsCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	articlePostIDsCachex = &ArticlePostIDsCachex{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
}

func (p *ArticlePostIDsCachex) GetRealDataFn() func(ctx context.Context, _ string) ([]int64, error) {
	return func(ctx context.Context, _ string) ([]int64, error) {
		db := mysql.GetDB(ctx)
		order, err := persistence.SelectPostOrderList(db)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return order, nil
	}
}

func (p *ArticlePostIDsCachex) Get(ctx context.Context) ([]int64, error) {
	order, ok := p.cacheX.Get(ctx, "", p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}
