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

// 初始化
var articleCachex *ArticleCachex

type ArticleCachex struct {
	cacheX *cachex.CacheX[int64, *dto.Article]
	expire time.Duration
}

func initArticleCachex() {
	cfg := config.GetCachexSettingByName("article")
	cx, err := infraCachex.NewCacheXBuilderByConfig[int64, *dto.Article](context.Background(), cfg).
		SetGetRealData((&ArticleCachex{}).GetRealData).
		SetMGetRealData((&ArticleCachex{}).MGetRealData).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	articleCachex = &ArticleCachex{
		expire: cfg.GetExpire(),
		cacheX: cx,
	}
}

func (a *ArticleCachex) GetRealData(ctx context.Context, id int64) (*dto.Article, error) {
	db := mysql.GetDB(ctx)
	// 获取post
	article, err := persistence.SelectArticleByID(db, id)
	if err != nil {
		return nil, infraCachex.ParseErr(err)
	}
	articleDTO := dto.NewArticleByModel(article)
	// 获取上一篇和下一篇
	articleDTO.PrevID, articleDTO.NextID = a.getPrevNext(ctx, id)
	// 获取作者
	user, _ := persistence.SelectUserByID(db, article.CreateUser)
	if user != nil {
		articleDTO.CreateUser = dto.NewUserByModel(user)
	}
	return articleDTO, nil
}

func (a *ArticleCachex) getPrevNext(ctx context.Context, id int64) (*int64, *int64) {
	list, err := articlePostIDsCachex.GetRealDataFn()(ctx, "")
	if err != nil {
		return nil, nil
	}
	var prev *int64
	var next *int64
	for i := 0; i < len(list); i++ {
		if i == len(list)-1 {
			next = nil
		} else {
			next = misc.ValPtr(list[i+1])
		}
		if list[i] == id {
			break
		}
		prev = misc.ValPtr(list[i])
	}
	return prev, next
}

func (a *ArticleCachex) MGetRealData(ctx context.Context, ids []int64) (map[int64]*dto.Article, error) {
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
	result := make(map[int64]*dto.Article)
	for id, article := range articles {
		articleDTO := dto.NewArticleByModel(article)
		articleDTO.PrevID, articleDTO.NextID = a.getPrevNext(ctx, article.ID)
		if user, ok := users[article.CreateUser]; ok {
			articleDTO.CreateUser = dto.NewUserByModel(user)
		}
		result[id] = articleDTO
	}
	return result, nil
}

func (a *ArticleCachex) MGet(ctx context.Context, ids []int64) map[int64]*dto.Article {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *ArticleCachex) Get(ctx context.Context, id int64) (*dto.Article, error) {
	article, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}

func (a *ArticleCachex) Refresh(ctx context.Context, id int64) {
	_ = a.cacheX.Delete(ctx, id)
	_, _ = a.Get(ctx, id)
}
