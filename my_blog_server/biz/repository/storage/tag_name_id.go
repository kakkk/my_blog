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

var tagNameIDStorage *TagNameIDStorage

type TagNameIDStorage struct {
	cacheX *cachex.CacheX[string, int64]
	expire time.Duration
}

func GetTagNameIDStorage() *TagNameIDStorage {
	return tagNameIDStorage
}

func initTagNameIDStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("tag_name_id")
	cache, err := NewCacheXBuilderByConfig[string, int64](ctx, cfg).
		SetGetRealData(tagNameIDGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	tagNameIDStorage = &TagNameIDStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func tagNameIDGetRealData(ctx context.Context, name string) (int64, error) {
	id, err := mysql.SelectTagIDByName(mysql.GetDB(ctx), name)
	if err != nil {
		return parseSqlError(id, err)
	}
	return id, nil
}

func (c *TagNameIDStorage) Get(ctx context.Context, name string) (int64, error) {
	id, ok := c.cacheX.Get(ctx, name, c.expire)
	if !ok {
		return id, consts.ErrRecordNotFound
	}
	return id, nil
}
