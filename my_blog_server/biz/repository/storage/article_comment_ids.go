package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/repository/mysql"
)

var articleCommentIDsStorage *ArticleCommentIDsStorage

type ArticleCommentIDsStorage struct {
	cacheX *cachex.CacheX[int64, []int64]
	expire time.Duration
}

func GetPostCommentIDsStorage() *ArticleCommentIDsStorage {
	return articleCommentIDsStorage
}

func initPostCommentIDsStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("article_comment_ids")
	cache, err := NewCacheXBuilderByConfig[int64, []int64](ctx, cfg).
		SetGetRealData(articleCommentIDsGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	articleCommentIDsStorage = &ArticleCommentIDsStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func articleCommentIDsGetRealData(ctx context.Context, id int64) ([]int64, error) {
	ids, err := mysql.SelectCommentIDsByArticleID(mysql.GetDB(ctx), id)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			// 缓存空值
			return make([]int64, 0), nil
		}
		return nil, err
	}
	return ids, err
}

func (a *ArticleCommentIDsStorage) Get(ctx context.Context, id int64) []int64 {
	ids, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return make([]int64, 0)
	}
	return ids
}
