package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/dto"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
)

var postPrevNextStorage *PostPrevNextStorage

type PostPrevNextStorage struct {
	cacheX *cachex.CacheX[int64, *dto.PostPrevNext]
	expire time.Duration
}

func GetPostPrevNextStorage() *PostPrevNextStorage {
	return postPrevNextStorage
}

func initPostPrevNextStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("post_pre_next")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, *dto.PostPrevNext](ctx, cfg).
		SetGetRealData(postPrevNextGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postPrevNextStorage = &PostPrevNextStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
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
	pn, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return pn, nil
}

// 重建缓存
func (p *PostPrevNextStorage) Rebuild(ctx context.Context, id int64) {
	_ = p.cacheX.Delete(ctx, id)
	_, _ = p.cacheX.Get(ctx, id, 0)
	return
}
