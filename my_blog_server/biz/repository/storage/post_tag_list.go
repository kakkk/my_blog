package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var postTagListStorage *PostTagListStorage

type PostTagListStorage struct {
	cacheX *cachex.CacheX[int64, []string]
}

func GetPostTagListStorage() *PostTagListStorage {
	return postTagListStorage
}

func initPostTagListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[[]string](ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache[[]string](ctx, 200, time.Minute)
	cache := cachex.NewCacheX[int64, []string]("post_tag_list", false, true).
		SetGetCacheKey(postTagListGetKey).
		SetGetRealData(postTagListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postTagListStorage = &PostTagListStorage{
		cacheX: cache,
	}
	return nil
}

func postTagListGetKey(id int64) string {
	return fmt.Sprintf("post_tag_list_%v", id)
}

func postTagListGetRealData(ctx context.Context, id int64) ([]string, error) {
	db := mysql.GetDB(ctx)
	tags, err := mysql.SelectTagListByArticleID(db, id)
	if err != nil {
		return parseSqlError(tags, err)
	}
	return tags, nil
}

func (p *PostTagListStorage) Get(ctx context.Context, id int64) ([]string, error) {
	tags, err := p.cacheX.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return tags, err
}

// 重建缓存
func (p *PostTagListStorage) Rebuild(ctx context.Context, id int64) {
	p.cacheX.Delete(ctx, id)
	_, _ = p.cacheX.Get(ctx, id)
	return
}
