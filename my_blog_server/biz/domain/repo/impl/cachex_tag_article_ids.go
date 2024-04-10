package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/config"
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

var tagArticleIDsCachex *TagArticleIDsCachex

// CategoryArticleIDsCachex key: tag_name, val: []article_id
type TagArticleIDsCachex struct {
	cacheX *cachex.CacheX[string, []int64]
	expire time.Duration
}

func initTagArticleIDsCachex() {
	cfg := config.GetCachexSettingByName("tag_article_ids")
	cache, err := infraCachex.NewCacheXBuilderByConfig[string, []int64](context.Background(), cfg).
		SetGetRealData((&TagArticleIDsCachex{}).GetRealDataFn()).
		Build()

	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	tagArticleIDsCachex = &TagArticleIDsCachex{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return
}

func (c *TagArticleIDsCachex) GetRealDataFn() func(ctx context.Context, name string) ([]int64, error) {
	return func(ctx context.Context, name string) ([]int64, error) {
		db := mysql.GetDB(ctx)
		tagID, err := persistence.SelectTagIDByName(db, name)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		list, err := persistence.SelectArticleIDsByTagIDs(mysql.GetDB(ctx), []int64{tagID})
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return list, nil
	}
}

func (c *TagArticleIDsCachex) Get(ctx context.Context, name string) ([]int64, error) {
	order, ok := c.cacheX.Get(ctx, name, c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}
