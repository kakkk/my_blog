package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/dto"
)

var articleMetaStorage *ArticleMetaStorage

type ArticleMetaStorage struct {
	cacheX *cachex.CacheX[int64, *dto.PostMeta]
	expire time.Duration
}

func GetArticleMetaStorage() *ArticleMetaStorage {
	return articleMetaStorage
}

func initArticleMetaStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("article_meta")
	cache, err := NewCacheXBuilderByConfig[int64, *dto.PostMeta](ctx, cfg).
		SetGetRealData(articleMetaStorageGetRealData).
		SetMGetRealData(articleMetaStorageMGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	articleMetaStorage = &ArticleMetaStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func articleMetaStorageGetRealData(ctx context.Context, id int64) (*dto.PostMeta, error) {
	post, err := GetArticleEntityStorage().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	editor, err := GetUserEntityStorage().Get(ctx, post.CreateUser)
	if err != nil {
		return nil, err
	}
	return dto.NewArticleMetaByEntity(post, editor), nil
}

func articleMetaStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*dto.PostMeta, error) {
	posts := GetArticleEntityStorage().MGet(ctx, ids)
	result := make(map[int64]*dto.PostMeta, len(posts))
	userIDs := make([]int64, 0, len(posts))
	for _, article := range posts {
		userIDs = append(userIDs, article.CreateUser)
	}
	users := GetUserEntityStorage().MGet(ctx, userIDs)
	for id, article := range posts {
		result[id] = dto.NewArticleMetaByEntity(article, users[article.ID])
	}
	return result, nil
}

func (p *ArticleMetaStorage) MGet(ctx context.Context, ids []int64) map[int64]*dto.PostMeta {
	return p.cacheX.MGet(ctx, ids, p.expire)
}

func (p *ArticleMetaStorage) Get(ctx context.Context, id int64) (*dto.PostMeta, error) {
	article, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
