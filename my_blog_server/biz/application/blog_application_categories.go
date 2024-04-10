package application

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/model/blog/page"
)

func (app *BlogApplication) CategoriesPage(ctx context.Context) (*page.TermsPageResp, *errorx.PageError) {
	categories, err := repo.GetContentRepo().Cache().GetCategoryList(ctx)
	if err != nil {
		return nil, errorx.NewInternalErrPageError()
	}
	if err != nil {
		log.GetLoggerWithCtx(ctx).Errorf("get category list error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	title := fmt.Sprintf("分类 - %v - %v", config.GetBlogName(), config.GetBlogSubTitle())
	return &page.TermsPageResp{
		List: dto.Categories(categories).ToTermListItem(),
		Meta: &page.PageMeta{
			Title:       title,
			Description: config.GetBlogDescription(),
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeCategoryList,
		},
	}, nil
}

func (app *BlogApplication) CategoryPostListPage(ctx context.Context, req *page.PostListPageRequest) (*page.PostListPageResp, *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"slug": req.GetName(),
		"page": req.GetPage(),
	})
	// 先拿分类
	category, err := repo.GetContentRepo().Cache().GetCategoryBySlug(ctx, req.GetName())
	if err != nil {
		logger.Warnf("get category fail: %v", err)
		// 不存在会返回 not found
		return nil, errorx.NewPageErrorByErr(err)
	}
	// 拿文章ID
	allArticleIDs, err := repo.GetContentRepo().Cache().GetArticleIDsByCategoryID(ctx, category.ID)
	if err != nil {
		logger.Warnf("get article ids fail: %v", err)
		return nil, errorx.NewPageErrorByErr(err)
	}
	// 拿当前页面的ID
	articleIDs, hasMore := misc.GetIDsByPage(allArticleIDs, req.GetPage())
	// 拿文章
	articlesMap := repo.GetContentRepo().Cache().MGetArticleMeta(ctx, articleIDs)
	// 排序
	articles := misc.MapToListByOrder(articleIDs, articlesMap)
	// 返回
	return dto.ArticleMetas(articles).PackPostListPageResp(req.GetPage(), hasMore, req.GetPageType(), category.CategoryName, category.Slug), nil
}
