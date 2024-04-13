package impl

import (
	"context"
	"fmt"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/cache"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/pkg/log"
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

func (c ContentCacheImpl) RefreshByPostID(ctx context.Context, id int64) {
	// 刷新article
	articleCachex.Refresh(ctx, id)
	// 刷新article meta
	articleMetaCachex.Refresh(ctx, id)
	// 刷新分类
	articleCategories.Refresh(ctx, id)
	// 刷新标签
	articleTagsCachex.Refresh(ctx, id)
	// 刷新分类下文章列表
	categories, err := articleCategories.Get(ctx, id)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("refresh cache get categories fail: %v", err)
	}
	categoryIDs := make([]int64, 0, len(categories))
	for _, category := range categories {
		categoryIDs = append(categoryIDs, category.ID)
	}
	categoryArticleIDsCachex.MRefresh(ctx, categoryIDs)
	// 刷新标签下文章列表
	tags, err := articleTagsCachex.Get(ctx, id)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("refresh cache get tags fail: %v", err)
	}
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.TagName)
	}
	tagArticleIDsCachex.MRefresh(ctx, tagNames)
	// 刷新文章列表
	articlePostIDsCachex.Refresh(ctx)
	// 设置UV
	_ = cache.GetPostUVCache().SetPostUVNX(ctx, redis.GetClient(), id, 0)
}

func (c ContentCacheImpl) RefreshByPageSlug(ctx context.Context, slug string) {
	articleSlugCachex.Refresh(ctx, slug)
}

func InitContentCache(ctx context.Context) {
	defaultID, err := persistence.SelectDefaultCategoryID(mysql.GetDB(ctx))
	if err != nil {
		panic(fmt.Errorf("select default category_id error:[%v]", err))
	}
	cache.SetDefaultCategoryID(defaultID)
}
