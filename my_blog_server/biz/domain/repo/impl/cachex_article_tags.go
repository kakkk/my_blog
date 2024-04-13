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

type ArticleTagsCachex struct {
	cacheX *cachex.CacheX[int64, []*dto.Tag]
	expire time.Duration
}

var articleTagsCachex *ArticleTagsCachex

func initArticleTagsCachex() {
	cfg := config.GetCachexSettingByName("article_tags")
	cx, err := infraCachex.NewCacheXBuilderByConfig[int64, []*dto.Tag](context.Background(), cfg).
		SetGetRealData((&ArticleTagsCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %v", err))
	}
	articleTagsCachex = &ArticleTagsCachex{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
}

func (a *ArticleTagsCachex) GetRealDataFn() func(ctx context.Context, id int64) ([]*dto.Tag, error) {
	return func(ctx context.Context, id int64) ([]*dto.Tag, error) {
		db := mysql.GetDB(ctx)
		ids, err := persistence.SelectTagIDsByArticleID(db, id)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		tags, err := persistence.MSelectTagByID(db, ids)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return dto.NewTagsByModelMap(tags), nil
	}
}

func (a *ArticleTagsCachex) Get(ctx context.Context, id int64) ([]*dto.Tag, error) {
	tags, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return tags, nil
}

func (a *ArticleTagsCachex) Refresh(ctx context.Context, id int64) {
	_ = a.cacheX.Delete(ctx, id)
	_, _ = a.cacheX.Get(ctx, id, a.expire)
}
