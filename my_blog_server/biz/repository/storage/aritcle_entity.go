package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/infra/config"
	cachex2 "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/model"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

var articleEntityStorage *ArticleEntityStorage

type ArticleEntityStorage struct {
	cacheX *cachex.CacheX[int64, *model.Article]
	expire time.Duration
}

func GetArticleEntityStorage() *ArticleEntityStorage {
	return articleEntityStorage
}

func initArticleEntityStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("article_entity")
	cx, err := cachex2.NewCacheXBuilderByConfig[int64, *model.Article](ctx, cfg).
		SetGetRealData(articleStorageGetRealData).
		SetMGetRealData(articleStorageMGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	articleEntityStorage = &ArticleEntityStorage{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
	return nil
}

func articleStorageGetRealData(ctx context.Context, id int64) (*model.Article, error) {
	db := mysql2.GetDB(ctx)
	// 获取post
	post, err := mysql.SelectArticleByID(db, id)
	if err != nil {
		return parseSqlError(post, err)
	}
	return post, nil
}

func articleStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*model.Article, error) {
	db := mysql2.GetDB(ctx)
	posts, err := mysql.MSelectPostWithPublishByIDs(db, ids)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, cachex.ErrNotFound
		}
		return nil, fmt.Errorf("sql error: %w", err)
	}
	return posts, nil
}

func (a *ArticleEntityStorage) MGet(ctx context.Context, ids []int64) map[int64]*model.Article {
	return a.cacheX.MGet(ctx, ids, a.expire)
}

func (a *ArticleEntityStorage) Get(ctx context.Context, id int64) (*model.Article, error) {
	article, ok := a.cacheX.Get(ctx, id, a.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return article, nil
}
