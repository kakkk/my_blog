package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
	"my_blog/biz/repository/dto"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/redis"
	"time"
)

type CategoryListStorage struct {
	cacheX *cachex.CacheX[*dto.CategoryList, int]
}

var categoryListStorage *CategoryListStorage

func initCategoryListStorage(ctx context.Context) error {
	redisCache := cachex.NewRedisCache(ctx, redis.GetRedisClient(ctx), 0)
	lruCache := cachex.NewLRUCache(ctx, 1, time.Hour)
	cache := cachex.NewCacheX[*dto.CategoryList, int]("category_list", false, false).
		SetGetCacheKey(categoryListGetKey).
		SetGetRealData(categoryListGetRealData).
		AddCache(ctx, true, lruCache).
		AddCache(ctx, false, redisCache)

	err := cache.Initialize(ctx)
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	categoryListStorage = &CategoryListStorage{
		cacheX: cache,
	}
	return nil
}

func categoryListGetKey(_ int) string {
	return "category_list"
}

func categoryListGetRealData(ctx context.Context, _ int) (*dto.CategoryList, error) {
	db := mysql.GetDB(ctx)
	order, err := mysql.SelectCategoryOrder(db)
	if err != nil {
		return nil, fmt.Errorf("select category order error:[%v]", err)
	}
	categories, err := mysql.MSelectCategoryByIDs(db, order)
	if err != nil {
		return nil, fmt.Errorf("select category error:[%v]", err)
	}

	counts, err := mysql.MSelectCategoryArticleCountByCategoryIDs(db, order)
	if err != nil {
		return nil, fmt.Errorf("select article count error:[%v]", err)
	}

	list := dto.CategoryList{}
	for _, id := range order {
		category, ok := categories[id]
		if !ok {
			logger.Warnf("category not exist, category_id:[%v]", id)
			continue
		}
		count := counts[id]
		list = append(list, &dto.CategoryListItem{
			ID:    category.ID,
			Name:  category.CategoryName,
			Slug:  category.Slug,
			Count: count,
		})
	}
	return &list, nil
}

func GetCategoryListStorage() *CategoryListStorage {
	return categoryListStorage
}

func (c *CategoryListStorage) Get(ctx context.Context) (*dto.CategoryList, error) {
	got, err := c.cacheX.Get(ctx, 0)
	if err != nil {
		if errors.Is(err, cachex.ErrNotFound) {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("get from cachex error:[%v]", err)
	}
	return got, nil
}

func (c *CategoryListStorage) GetFromDB(ctx context.Context) (*dto.CategoryList, error) {
	return categoryListGetRealData(ctx, 0)
}

func (c *CategoryListStorage) RebuildCache(ctx context.Context) {
	c.cacheX.Delete(ctx, 0)
	_, _ = c.cacheX.Get(ctx, 0)
}
