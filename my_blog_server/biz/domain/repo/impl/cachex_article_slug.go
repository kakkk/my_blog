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

// 初始化
var articleSlugCachex *ArticleSlugCachex

type ArticleSlugCachex struct {
	cacheX *cachex.CacheX[string, *dto.Article]
	expire time.Duration
}

func initArticleSlugCachex() {
	cfg := config.GetCachexSettingByName("article_slug")
	cx, err := infraCachex.NewCacheXBuilderByConfig[string, *dto.Article](context.Background(), cfg).
		SetGetRealData((&ArticleSlugCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	articleSlugCachex = &ArticleSlugCachex{
		expire: cfg.GetExpire(),
		cacheX: cx,
	}
}

func (a *ArticleSlugCachex) GetRealDataFn() func(ctx context.Context, slug string) (*dto.Article, error) {
	return func(ctx context.Context, slug string) (*dto.Article, error) {
		db := mysql.GetDB(ctx)
		// 获取post
		article, err := persistence.SelectArticleBySlug(db, slug)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		articleDTO := dto.NewArticleByModel(article)
		// 获取作者
		user, _ := persistence.SelectUserByID(db, article.CreateUser)
		if user != nil {
			articleDTO.CreateUser = dto.NewUserByModel(user)
		}
		return articleDTO, nil
	}
}

func (a *ArticleSlugCachex) Get(ctx context.Context, slug string) (*dto.Article, error) {
	article, ok := a.cacheX.Get(ctx, slug, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}

func (a *ArticleSlugCachex) Refresh(ctx context.Context, slug string) {
	_ = a.cacheX.Delete(ctx, slug)
	_, _ = a.Get(ctx, slug)
}
