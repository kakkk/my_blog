package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

type ArticleCommentsCachex struct {
	cacheX *cachex.CacheX[int64, []*dto.Comment]
	expire time.Duration
}

var articleCommentsCachex *ArticleCommentsCachex

func initArticleCommentsCachex() {
	cfg := config.GetCachexSettingByName("article_comments")
	cx, err := infraCachex.NewCacheXBuilderByConfig[int64, []*dto.Comment](context.Background(), cfg).
		SetGetRealData((&ArticleCommentsCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	articleCommentsCachex = &ArticleCommentsCachex{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
}

func (a *ArticleCommentsCachex) GetRealDataFn() func(ctx context.Context, id int64) ([]*dto.Comment, error) {
	return func(ctx context.Context, id int64) ([]*dto.Comment, error) {
		db := mysql.GetDB(ctx)
		commentIDs, err := persistence.SelectCommentIDsByArticleID(db, id)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		comments := commentCachex.MGet(ctx, commentIDs)
		return misc.MapValues(comments), nil
	}
}

func (a *ArticleCommentsCachex) Get(ctx context.Context, id int64) []*dto.Comment {
	comments, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return make([]*dto.Comment, 0)
	}
	return comments
}

func (a *ArticleCommentsCachex) Refresh(ctx context.Context, id int64) {
	_ = a.cacheX.Delete(ctx, id)
	_ = a.Get(ctx, id)
}
