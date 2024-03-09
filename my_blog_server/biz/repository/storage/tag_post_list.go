package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

var tagPostListStorage *TagPostListStorage

type TagPostListStorage struct {
	cacheX *cachex.CacheX[int64, []int64]
	expire time.Duration
}

func GetTagPostListStorage() *TagPostListStorage {
	return tagPostListStorage
}

func initTagPostListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("tag_post_list")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, []int64](ctx, cfg).
		SetGetRealData(tagPostListGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	tagPostListStorage = &TagPostListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func tagPostListGetRealData(ctx context.Context, id int64) ([]int64, error) {
	list, err := mysql.SelectPostIDsByTagID(mysql2.GetDB(ctx), id)
	if err != nil {
		return parseSqlError(list, err)
	}
	return list, nil
}

func (p *TagPostListStorage) Get(ctx context.Context, id int64) ([]int64, error) {
	order, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}

// 重建缓存
func (p *TagPostListStorage) Rebuild(ctx context.Context, ids []int64) {
	for _, id := range ids {
		_ = p.cacheX.Delete(ctx, id)
		_, _ = p.cacheX.Get(ctx, id, 0)
	}
	return
}
