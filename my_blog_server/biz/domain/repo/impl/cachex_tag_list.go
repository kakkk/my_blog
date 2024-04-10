package impl

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/config"
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

type TagListCachex struct {
	cacheX *cachex.CacheX[string, []*dto.Tag]
	expire time.Duration
}

var tagListCachex *TagListCachex

func initTagListCachex() {
	cfg := config.GetCachexSettingByName("tag_list")
	cache, err := infraCachex.NewCacheXBuilderByConfig[string, []*dto.Tag](context.Background(), cfg).
		SetGetRealData((&TagListCachex{}).GetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	tagListCachex = &TagListCachex{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
}

func (*TagListCachex) GetRealDataFn() func(ctx context.Context, _ string) ([]*dto.Tag, error) {
	return func(ctx context.Context, _ string) ([]*dto.Tag, error) {
		db := mysql.GetDB(ctx)
		tagModelList, err := persistence.GetAllTag(db)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		ids := make([]int64, 0, len(tagModelList))
		for _, tag := range tagModelList {
			ids = append(ids, tag.ID)
		}
		tagPostCountMap, err := persistence.MGetTagArticleCountByTagIDs(db, ids, true)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		result := make([]*dto.Tag, 0, len(tagModelList))
		for _, tag := range tagModelList {
			result = append(result, &dto.Tag{
				ID:      tag.ID,
				TagName: tag.TagName,
				Count:   tagPostCountMap[tag.ID],
			})
		}
		// 排序
		sort.Slice(result, func(i, j int) bool {
			return result[i].Count >= result[j].Count
		})
		return result, nil
	}
}

func (c *TagListCachex) Get(ctx context.Context) ([]*dto.Tag, error) {
	got, ok := c.cacheX.Get(ctx, "", c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return got, nil
}

func (c *TagListCachex) RebuildCache(ctx context.Context) {
	_ = c.cacheX.Delete(ctx, "")
	_, _ = c.cacheX.Get(ctx, "", 0)
}
