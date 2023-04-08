package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
)

var postOrderListStorage *postOrderList

type orderList []int64

func (a *orderList) Serialize() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a *orderList) Deserialize(str string) (*orderList, error) {
	list := &orderList{}
	err := json.Unmarshal([]byte(str), list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

type postOrderList struct {
	cacheX *cachex.CacheX[*orderList, int]
}

func GetPostOrderListStorage() *postOrderList {
	return postOrderListStorage
}

func initPostOrderListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), time.Minute*30)
	lruCache := cachex.NewLRUCache(ctx, 1, time.Minute)
	cache := cachex.NewCacheX[*orderList, int]("post_order_list", false, true).
		SetGetCacheKey(postOrderListGetKey).
		SetGetRealData(postOrderListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postOrderListStorage = &postOrderList{
		cacheX: cache,
	}
	return nil
}

func postOrderListGetKey(_ int) string {
	return "post_order_list"
}

func postOrderListGetRealData(ctx context.Context, _ int) (*orderList, error) {
	db := mysql.GetDB(ctx)
	order, err := mysql.SelectPostOrderList(db)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("sql error:[%w]", err)
	}
	return (*orderList)(&order), nil
}

func (p *postOrderList) Get(ctx context.Context) ([]int64, error) {
	order, err := p.cacheX.Get(ctx, 0)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return *order, err
}

// 重建缓存
func (p *postOrderList) Rebuild(ctx context.Context) {
	p.cacheX.Delete(ctx, 0)
	_, _ = p.cacheX.Get(ctx, 0)
	return
}
