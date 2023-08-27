package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/repository/mysql"
)

var postOrderListStorage *PostOrderListStorage

type PostOrderListStorage struct {
	cacheX *cachex.CacheX[string, []int64]
	expire time.Duration
}

func GetPostOrderListStorage() *PostOrderListStorage {
	return postOrderListStorage
}

func initPostOrderListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("post_order_list")
	cache, err := NewCacheXBuilderByConfig[string, []int64](ctx, cfg).
		SetGetRealData(postOrderListGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postOrderListStorage = &PostOrderListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func postOrderListGetRealData(ctx context.Context, _ string) ([]int64, error) {
	db := mysql.GetDB(ctx)
	order, err := mysql.SelectPostOrderList(db)
	if err != nil {
		return parseSqlError(order, err)
	}
	return order, nil
}

func (p *PostOrderListStorage) Get(ctx context.Context) ([]int64, error) {
	order, ok := p.cacheX.Get(ctx, "", p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}

// 重建缓存
func (p *PostOrderListStorage) Rebuild(ctx context.Context) {
	_ = p.cacheX.Delete(ctx, "")
	_, _ = p.cacheX.Get(ctx, "", 0)
	return
}
