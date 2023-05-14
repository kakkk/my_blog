package storage

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/components/cachex"
	"my_blog/biz/dto"
)

var postPrevNextStorage *PostPrevNextStorage

type PostPrevNextStorage struct {
	cacheX *cachex.CacheX[int64, *dto.PostPrevNext]
}

func GetPostPrevNextStorage() *PostPrevNextStorage {
	return postPrevNextStorage
}

func initPostPrevNextStorage(ctx context.Context) error {
	lruCache := cachex.NewLRUCache[*dto.PostPrevNext](ctx, 200, 5*time.Minute)
	cache := cachex.NewCacheX[int64, *dto.PostPrevNext]("post_order_list", false, true).
		SetGetCacheKey(postPrevNextGetKey).
		SetGetRealData(postPrevNextGetRealData).
		AddCache(ctx, true, lruCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postPrevNextStorage = &PostPrevNextStorage{
		cacheX: cache,
	}
	return nil
}

func postPrevNextGetKey(id int64) string {
	return fmt.Sprintf("post_prev_next_%v", id)
}

func postPrevNextGetRealData(ctx context.Context, id int64) (*dto.PostPrevNext, error) {
	list, err := GetPostOrderListStorage().Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get order list error")
	}
	var prev *int64
	var next *int64
	for i := 0; i < len(list); i++ {
		if i == len(list)-1 {
			next = nil
		} else {
			next = &list[i+1]
		}
		if list[i] == id {
			break
		}
		prev = &list[i]
	}
	return &dto.PostPrevNext{Prev: prev, Next: next}, nil
}

func (p *PostPrevNextStorage) Get(ctx context.Context, id int64) (*dto.PostPrevNext, error) {
	pn, err := p.cacheX.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get from cachex error:[%w]", err)
	}
	return pn, err
}

// 重建缓存
func (p *PostPrevNextStorage) Rebuild(ctx context.Context, id int64) {
	p.cacheX.Delete(ctx, id)
	_, _ = p.cacheX.Get(ctx, id)
	return
}
