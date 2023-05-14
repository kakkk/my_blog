package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/dto"
	"my_blog/biz/repository/redis"
)

var postMetaStorage *PostMetaStorage

type PostMetaStorage struct {
	cacheX *cachex.CacheX[int64, *dto.PostMeta]
}

func GetPostMetaStorage() *PostMetaStorage {
	return postMetaStorage
}

func initPostMetaStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[*dto.PostMeta](ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache[*dto.PostMeta](ctx, 1000, time.Minute)
	cache := cachex.NewCacheX[int64, *dto.PostMeta]("article_entity", false, true).
		SetGetCacheKey(postMetaStorageGetKey).
		SetGetRealData(postMetaStorageGetRealData).
		SetMGetRealData(postMetaStorageMGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postMetaStorage = &PostMetaStorage{
		cacheX: cache,
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

func postMetaStorageGetKey(id int64) string {
	return fmt.Sprintf("my_blog_post_meta_%v", id)
}

func (p *PostMetaStorage) MGet(ctx context.Context, ids []int64) map[int64]*dto.PostMeta {
	return p.cacheX.MGet(ctx, ids)
}

func (p *PostMetaStorage) Get(ctx context.Context, id int64) (*dto.PostMeta, error) {
	article, err := p.cacheX.Get(ctx, id)
	if errors.Is(err, cachex.ErrNotFound) {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
