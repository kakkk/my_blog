package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/entity"
	"my_blog/biz/repository/mysql"
)

var postCategoryListStorage *PostCategoryListStorage

type PostCategoryListStorage struct {
	cacheX *cachex.CacheX[int64, []*entity.Category]
	expire time.Duration
}

func GetPostCategoryListStorage() *PostCategoryListStorage {
	return postCategoryListStorage
}

func initPostCategoryListStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("post_category_list")
	cache, err := NewCacheXBuilderByConfig[int64, []*entity.Category](ctx, cfg).
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

func postCategoryListGetRealData(ctx context.Context, id int64) ([]*entity.Category, error) {
	db := mysql.GetDB(ctx)
	categoryIDs, err := mysql.SelectCategoryIDsByArticleID(db, id)
	if err != nil {
		return parseSqlError([]*entity.Category{}, err)
	}
	cMap := GetCategoryEntityStorage().MGet(ctx, categoryIDs)
	result := make([]*entity.Category, 0, len(cMap))
	for _, category := range cMap {
		result = append(result, category)
	}
	return result, nil
}

func (p *PostCategoryListStorage) Get(ctx context.Context, id int64) ([]*entity.Category, error) {
	list, ok := p.cacheX.Get(ctx, id, p.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return list, nil
}
