package storage

import (
	"context"
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

var postCategoryListStorage *PostCategoryListStorage

type PostCategoryListStorage struct {
	cacheX *cachex.CacheX[int64, []*model.Category]
	expire time.Duration
}

func GetPostCategoryListStorage() *PostCategoryListStorage {
	return postCategoryListStorage
}

func initPostCategoryListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("post_category_list")
	cache, err := cachex2.NewCacheXBuilderByConfig[int64, []*model.Category](ctx, cfg).
		SetGetRealData(postCategoryListGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	postCategoryListStorage = &PostCategoryListStorage{
		cacheX: cache,
		expire: cfg.GetExpire(),
	}
	return nil
}

func postCategoryListGetRealData(ctx context.Context, id int64) ([]*model.Category, error) {
	db := mysql2.GetDB(ctx)
	categoryIDs, err := mysql.SelectCategoryIDsByArticleID(db, id)
	if err != nil {
		return parseSqlError([]*model.Category{}, err)
	}
	cMap := GetCategoryEntityStorage().MGet(ctx, categoryIDs)
	result := make([]*model.Category, 0, len(cMap))
	for _, category := range cMap {
		result = append(result, category)
	}
	return result, nil
}

func (p *PostCategoryListStorage) Get(ctx context.Context, id int64) ([]*model.Category, error) {
	list, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return list, nil
}
