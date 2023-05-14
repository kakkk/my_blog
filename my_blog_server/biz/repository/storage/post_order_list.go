package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var postOrderListStorage *PostOrderListStorage

type PostOrderListStorage struct {
	cacheX *cachex.CacheX[int, []int64]
}

func GetPostOrderListStorage() *PostOrderListStorage {
	return postOrderListStorage
}

func initPostOrderListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache[[]int64](ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache[[]int64](ctx, 1, time.Minute)
	cache := cachex.NewCacheX[int, []int64]("post_order_list", false, true).
		SetGetCacheKey(postOrderListGetKey).
		SetGetRealData(postOrderListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postOrderListStorage = &PostOrderListStorage{
		cacheX: cache,
	}
	return nil
}

func postOrderListGetKey(_ int) string {
	return "post_order_list"
}

func postOrderListGetRealData(ctx context.Context, _ int) ([]int64, error) {
	db := mysql.GetDB(ctx)
	order, err := mysql.SelectPostOrderList(db)
	if err != nil {
		return parseSqlError(order, err)
	}
	return order, nil
}

func (p *PostOrderListStorage) Get(ctx context.Context) ([]int64, error) {
	order, err := p.cacheX.Get(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return order, err
}

// 重建缓存
func (p *PostOrderListStorage) Rebuild(ctx context.Context) {
	p.cacheX.Delete(ctx, 0)
	_, _ = p.cacheX.Get(ctx, 0)
	return
}
