package application

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/hertz_gen/blog/page"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
)

func (app *BlogApplication) TagsPage(ctx context.Context) (rsp *page.TermsPageResp, pErr *errorx.PageError) {
	tags, err := repo.GetContentRepo().Cache().GetTagList(ctx)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Errorf("get tag list error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	title := fmt.Sprintf("标签 - %v - %v", config.GetBlogName(), config.GetBlogSubTitle())
	return &page.TermsPageResp{
		List: (dto.Tags(tags)).ToTermListItem(),
		Meta: &page.PageMeta{
			Title:       title,
			Description: config.GetBlogDescription(),
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeTagList,
		},
	}, nil
}

func (app *BlogApplication) TagPostListPage(ctx context.Context, req *page.PostListPageRequest) (rsp *page.PostListPageResp, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"slug": req.GetName(),
		"page": req.GetPage(),
	})
	// 拿文章ID
	allArticleIDs, err := repo.GetContentRepo().Cache().GetArticleIDsByTagName(ctx, req.GetName())
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
	return dto.ArticleMetas(articles).PackPostListPageResp(req.GetPage(), hasMore, req.GetPageType(), req.GetName(), req.GetName()), nil
}
