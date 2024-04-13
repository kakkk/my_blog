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
	"my_blog/biz/infra/misc"
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

type ArticleMetaCachex struct {
	cacheX *cachex.CacheX[int64, *dto.ArticleMeta]
	expire time.Duration
}

var articleMetaCachex *ArticleMetaCachex

func initArticleMetaCachex() {
	cfg := config.GetCachexSettingByName("article_meta")
	cx, err := infraCachex.NewCacheXBuilderByConfig[int64, *dto.ArticleMeta](context.Background(), cfg).
		SetGetRealData((&ArticleMetaCachex{}).GetRealData).
		SetMGetRealData((&ArticleMetaCachex{}).MGetRealData).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	articleMetaCachex = &ArticleMetaCachex{
		expire: cfg.GetExpire(),
		cacheX: cx,
	}
}

func (a *ArticleMetaCachex) GetRealData(ctx context.Context, id int64) (*dto.ArticleMeta, error) {
	db := mysql.GetDB(ctx)
	// 获取post
	article, err := persistence.SelectArticleByID(db, id)
	if err != nil {
		return nil, infraCachex.ParseErr(err)
	}
	articleDTO := dto.NewArticleByModel(article)
	// 获取作者
	user, _ := persistence.SelectUserByID(db, article.CreateUser)
	if user != nil {
		articleDTO.CreateUser = dto.NewUserByModel(user)
	}
	return articleDTO.ToArticleMeta(), nil
}

func (a *ArticleMetaCachex) MGetRealData(ctx context.Context, ids []int64) (map[int64]*dto.ArticleMeta, error) {
	db := mysql.GetDB(ctx)
	articles, err := persistence.MSelectPostWithPublishByIDs(db, ids)
	if err != nil {
		return nil, infraCachex.ParseErr(err)
	}
	userIDs := make(map[int64]struct{})
	for _, article := range articles {
		userIDs[article.CreateUser] = struct{}{}
	}
	users, _ := persistence.MSelectUserByIDs(db, misc.MapKeys(userIDs))
	result := make(map[int64]*dto.ArticleMeta)
	for id, article := range articles {
		articleDTO := dto.NewArticleByModel(article)
		if user, ok := users[article.CreateUser]; ok {
			articleDTO.CreateUser = dto.NewUserByModel(user)
		}
		result[id] = articleDTO.ToArticleMeta()
	}
	return result, nil
}

func (a *ArticleMetaCachex) MGet(ctx context.Context, ids []int64) map[int64]*dto.ArticleMeta {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *ArticleMetaCachex) Get(ctx context.Context, id int64) (*dto.ArticleMeta, error) {
	article, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}

func (a *ArticleMetaCachex) Refresh(ctx context.Context, id int64) {
	_ = a.cacheX.Delete(ctx, id)
	_, _ = a.Get(ctx, id)
}
