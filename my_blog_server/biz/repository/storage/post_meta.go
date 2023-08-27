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

var postMetaStorage *PostMetaStorage

type PostMetaStorage struct {
	cacheX *cachex.CacheX[int64, *dto.PostMeta]
	expire time.Duration
}

func GetPostMetaStorage() *PostMetaStorage {
	return postMetaStorage
}

func initPostMetaStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("post_meta")
	cache, err := NewCacheXBuilderByConfig[int64, *dto.PostMeta](ctx, cfg).
		SetGetRealData(postMetaStorageGetRealData).
		SetMGetRealData(postMetaStorageMGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postMetaStorage = &PostMetaStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func postMetaStorageGetRealData(ctx context.Context, id int64) (*dto.PostMeta, error) {
	post, err := GetArticleEntityStorage().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	editor, err := GetUserEntityStorage().Get(ctx, post.CreateUser)
	if err != nil {
		return nil, err
	}
	return dto.NewPostMetaByEntity(post, editor), nil
}

func postMetaStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*dto.PostMeta, error) {
	posts := GetArticleEntityStorage().MGet(ctx, ids)
	result := make(map[int64]*dto.PostMeta, len(posts))
	userIDs := make([]int64, 0, len(posts))
	for _, article := range posts {
		userIDs = append(userIDs, article.CreateUser)
	}
	users := GetUserEntityStorage().MGet(ctx, userIDs)
	for id, article := range posts {
		result[id] = dto.NewPostMetaByEntity(article, users[article.ID])
	}
	return result, nil
}

func (p *PostMetaStorage) MGet(ctx context.Context, ids []int64) map[int64]*dto.PostMeta {
	return p.cacheX.MGet(ctx, ids, p.expire)
}

func (p *PostMetaStorage) Get(ctx context.Context, id int64) (*dto.PostMeta, error) {
	article, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
