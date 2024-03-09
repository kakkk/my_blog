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

var postTagListStorage *PostTagListStorage

type PostTagListStorage struct {
	cacheX *cachex.CacheX[int64, []string]
	expire time.Duration
}

func GetPostTagListStorage() *PostTagListStorage {
	return postTagListStorage
}

func initPostTagListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("post_tag_list")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, []string](ctx, cfg).
		SetGetRealData(postTagListGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postTagListStorage = &PostTagListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func postTagListGetRealData(ctx context.Context, id int64) ([]string, error) {
	db := mysql2.GetDB(ctx)
	tags, err := mysql.SelectTagListByArticleID(db, id)
	if err != nil {
		return parseSqlError(tags, err)
	}
	return tags, nil
}

func (p *PostTagListStorage) Get(ctx context.Context, id int64) ([]string, error) {
	tags, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return tags, nil
}

// 重建缓存
func (p *PostTagListStorage) Rebuild(ctx context.Context, id int64) {
	_ = p.cacheX.Delete(ctx, id)
	_, _ = p.cacheX.Get(ctx, id, 0)
	return
}
