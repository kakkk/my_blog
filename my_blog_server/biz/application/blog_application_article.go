package application

import (
	"context"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/model/blog/page"
)

func (app *BlogApplication) PostPage(ctx context.Context, req *page.PostPageRequest) (rsp *page.PostPageResponse, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx).WithField("article_id", req.GetID())
	rsp = page.NewPostPageResponse()
	defer misc.Recover(ctx, func() {
		pErr = errorx.NewInternalErrPageError()
		return
	})()
	article := entity.NewArticlePageByID(req.GetID())
	// 获取内容
	err := article.FetchContent(ctx)
	if err != nil {
		logger.Warnf("fetch content fail: %v", err)
		return nil, errorx.NewPageErrorByErr(err)
	}
	// 处理UV
	article.View(ctx)
	// 打包
	rsp = article.PackPostPageResponse()
	return rsp, nil
}

func (app *BlogApplication) PostListPage(ctx context.Context, req *page.PostListPageRequest) (rsp *page.PostListPageResp, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx)
	allPostIDs, err := repo.GetContentRepo().Cache().GetArticlePostIDs(ctx)
	if err != nil {
		logger.Errorf("get post list error:[%v]", err)
		return nil, errorx.NewPageErrorByErr(err)
	}
	articleIDs, hasMore := misc.GetIDsByPage(allPostIDs, req.GetPage())
	if len(articleIDs) == 0 {
		return nil, errorx.NewNotFoundErrPageError()
	}
	articlesMap := repo.GetContentRepo().Cache().MGetArticleMeta(ctx, articleIDs)
	if len(articlesMap) == 0 {
		return nil, errorx.NewNotFoundErrPageError()
	}
	articles := misc.MapToListByOrder(articleIDs, articlesMap)
	return dto.ArticleMetas(articles).PackPostListPageResp(req.GetPage(), hasMore, req.GetPageType(), "", ""), nil
}
