package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/dto"
	"my_blog/biz/repository/mysql"
)

type TagListStorage struct {
	cacheX *cachex.CacheX[string, *dto.TagList]
	expire time.Duration
}

var tagListStorage *TagListStorage

func initTagListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("tag_list")
	cache, err := NewCacheXBuilderByConfig[string, *dto.TagList](ctx, cfg).
		SetGetRealData(tagListGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	tagListStorage = &TagListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func GetTagListStorage() *TagListStorage {
	return tagListStorage
}

func tagListGetRealData(ctx context.Context, _ string) (*dto.TagList, error) {
	db := mysql.GetDB(ctx)
	tagEntityList, err := mysql.GetAllTag(db)
	if err != nil {
		return parseSqlError(&dto.TagList{}, err)
	}
	ids := make([]int64, 0, len(tagEntityList))
	for _, tag := range tagEntityList {
		ids = append(ids, tag.ID)
	}
	tagPostCountMap, err := mysql.MGetTagArticleCountByTagIDs(db, ids, true)
	if err != nil {
		return parseSqlError(&dto.TagList{}, err)
	}
	return dto.NewTagList(tagEntityList, tagPostCountMap), nil
}

func (c *TagListStorage) Get(ctx context.Context) (*dto.TagList, error) {
	got, ok := c.cacheX.Get(ctx, "", c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return got, nil
}

func (c *TagListStorage) RebuildCache(ctx context.Context) {
	_ = c.cacheX.Delete(ctx, "")
	_, _ = c.cacheX.Get(ctx, "", 0)
}
