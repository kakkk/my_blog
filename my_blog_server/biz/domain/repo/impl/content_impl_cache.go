package impl

import (
	"context"
	"fmt"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/cache"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/infra/repository/redis"
)

type ContentCacheImpl struct{}

func (c ContentCacheImpl) GetCategoryList(ctx context.Context) ([]*dto.Category, error) {
	return categoryListCachex.Get(ctx)
}

func (c ContentCacheImpl) GetDefaultCategoryID() int64 {
	return cache.GetDefaultCategoryID()
}

func (c ContentCacheImpl) GetPostUV(ctx context.Context, id int64) (int64, error) {
	return cache.GetPostUVCache().Get(ctx, redis.GetClient(), id)
}

func (c ContentCacheImpl) IncrPostUV(ctx context.Context, id int64) error {
	return cache.GetPostUVCache().Incr(ctx, redis.GetClient(), id)
}

func (c ContentCacheImpl) GetArticle(ctx context.Context, id int64) (*dto.Article, error) {
	return articleCachex.Get(ctx, id)
}

func (c ContentCacheImpl) GetArticleBySlug(ctx context.Context, slug string) (*dto.Article, error) {
	return articleSlugCachex.Get(ctx, slug)
}

func (c ContentCacheImpl) GetArticleMeta(ctx context.Context, id int64) (*dto.ArticleMeta, error) {
	return articleMetaCachex.Get(ctx, id)
}

func (c ContentCacheImpl) GetArticlePostIDs(ctx context.Context) ([]int64, error) {
	return articlePostIDsCachex.Get(ctx)
}

func (c ContentCacheImpl) MGetArticleMeta(ctx context.Context, ids []int64) map[int64]*dto.ArticleMeta {
	return articleMetaCachex.MGet(ctx, ids)
}

func (c ContentCacheImpl) GetTagList(ctx context.Context) ([]*dto.Tag, error) {
	return tagListCachex.Get(ctx)
}

func (c ContentCacheImpl) GetCategoriesByArticleID(ctx context.Context, id int64) ([]*dto.Category, error) {
	return articleCategories.Get(ctx, id)
}

func (c ContentCacheImpl) GetTagsByArticleID(ctx context.Context, id int64) ([]*dto.Tag, error) {
	return articleTagsCachex.Get(ctx, id)
}

func (c ContentCacheImpl) GetCategoryBySlug(ctx context.Context, slug string) (*dto.Category, error) {
	return categoryCachex.Get(ctx, slug)
}

func (c ContentCacheImpl) GetArticleIDsByCategoryID(ctx context.Context, id int64) ([]int64, error) {
	return categoryArticleIDsCachex.Get(ctx, id)
}

func (c ContentCacheImpl) GetArticleIDsByTagName(ctx context.Context, name string) ([]int64, error) {
	return tagArticleIDsCachex.Get(ctx, name)
}

func InitContentCache(ctx context.Context) {
	defaultID, err := persistence.SelectDefaultCategoryID(mysql.GetDB(ctx))
	if err != nil {
		panic(fmt.Errorf("select default category_id error:[%v]", err))
	}
	cache.SetDefaultCategoryID(defaultID)
}
