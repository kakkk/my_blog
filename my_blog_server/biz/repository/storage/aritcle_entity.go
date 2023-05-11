package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var articleEntityStorage *ArticleEntityStorage

type ArticleEntityStorage struct {
	cacheX *cachex.CacheX[*entity.Article, int64]
}

func GetArticleEntityStorage() *ArticleEntityStorage {
	return articleEntityStorage
}

func initArticleEntityStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1000, time.Minute)
	cache := cachex.NewSerializableCacheX[*entity.Article, int64]("article_entity", false, true).
		SetGetCacheKey(articleStorageGetKey).
		SetGetRealData(articleStorageGetRealData).
		SetMGetRealData(articleStorageMGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	articleEntityStorage = &ArticleEntityStorage{
		cacheX: cache,
	}
	return nil
}

func articleStorageGetRealData(ctx context.Context, id int64) (*entity.Article, error) {
	db := mysql.GetDB(ctx)
	// 获取post
	post, err := mysql.SelectArticleByID(db, id)
	if err != nil {
		return parseSqlError(post, err)
	}
	return post, nil
}

func articleStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*entity.Article, error) {
	db := mysql.GetDB(ctx)
	posts, err := mysql.MSelectPostWithPublishByIDs(db, ids)
	if err != nil {
		if err == consts.ErrRecordNotFound {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return posts, nil
}

func articleStorageGetKey(id int64) string {
	return fmt.Sprintf("my_blog_article_entity_%v", id)
}

func (a *ArticleEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*entity.Article {
	return a.cacheX.MGet(ctx, ids)
}

func (a *ArticleEntityStorage) Get(ctx context.Context, id int64) (*entity.Article, error) {
	article, err := a.cacheX.Get(ctx, id)
	if err != nil {
		if errors.Is(err, cachex.ErrNotFound) {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("cacheX error:[%w]", err)
	}

	return article, nil
}
