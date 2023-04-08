package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/dto"
	"my_blog/biz/repository/redis"
)

var postMetaStorage *postMeta

type postMeta struct {
	cacheX *cachex.CacheX[*dto.PostMeta, int64]
}

func GetPostMetaStorage() *postMeta {
	return postMetaStorage
}

func initPostMetaStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1000, time.Minute)
	cache := cachex.NewCacheX[*dto.PostMeta, int64]("article_entity", false, true).
		SetGetCacheKey(postMetaStorageGetKey).
		SetGetRealData(postMetaStorageGetRealData).
		SetMGetRealData(postMetaStorageMGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postMetaStorage = &postMeta{
		cacheX: cache,
	}
	return nil
}

func postMetaStorageGetRealData(ctx context.Context, id int64) (*dto.PostMeta, error) {
	post, err := GetArticleEntityStorage().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.PostMeta{
		ID:    post.ID,
		Title: post.Title,
	}, nil
}

func postMetaStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*dto.PostMeta, error) {
	posts := GetArticleEntityStorage().MGet(ctx, ids)
	result := make(map[int64]*dto.PostMeta, len(posts))
	for id, article := range posts {
		result[id] = &dto.PostMeta{
			ID:    article.ID,
			Title: article.Title,
		}
	}
	return result, nil
}

func postMetaStorageGetKey(id int64) string {
	return fmt.Sprintf("my_blog_post_meta_%v", id)
}

func (p *postMeta) MGet(ctx context.Context, ids []int64) map[int64]*dto.PostMeta {
	return p.cacheX.MGet(ctx, ids)
}

func (p *postMeta) Get(ctx context.Context, id int64) (*dto.PostMeta, error) {
	article, err := p.cacheX.Get(ctx, id)
	if errors.Is(err, cachex.ErrNotFound) {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
