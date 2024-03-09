package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"my_blog/biz/common/resp"
	"my_blog/biz/consts"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/storage"
)

func CategoryPostListByPage(ctx context.Context, req *page.PostListPageRequest) (rsp *page.PostListPageResp, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"slug": req.GetName(),
		"page": req.GetPage(),
	})
	rsp = &page.PostListPageResp{
		Meta:     resp.NewBasePageMeta(page.PageTypeCategoryPostList),
		PostList: []*page.PostItem{},
	}
	defer misc.Recover(ctx, func() {
		pErr = errorx.NewInternalErrPageError()
		return
	})()
	// 获取分类id
	categoryID, err := storage.GetCategorySlugIDStorage().Get(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("slug not found")
			return nil, errorx.NewNotFoundErrPageError()
		}
		logger.Errorf("get category_id error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	// 获取分类
	category, err := storage.GetCategoryEntityStorage().Get(ctx, categoryID)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("slug not found")
			return nil, errorx.NewNotFoundErrPageError()
		}
		logger.Errorf("get category_entity error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	rsp.Meta.Title = fmt.Sprintf("分类 %v 下的文章 - %v", category.CategoryName, config.GetBlogName())
	rsp.Name = category.CategoryName
	rsp.Slug = category.Slug
	// 获取分类下所有文章id
	postIDs, err := storage.GetCategoryPostListStorage().Get(ctx, categoryID)
	if err != nil {
		if !errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("no article")
			return rsp, nil
		}
		logger.Errorf("get category_post_list error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	postIDs, hasMore := getIDsByPage(postIDs, req.GetPage())
	if len(postIDs) == 0 {
		logger.Warnf("no article")
		return rsp, nil
	}
	postMetas := storage.GetArticleMetaStorage().MGet(ctx, postIDs)
	if len(postMetas) == 0 {
		logger.Warnf("no article")
		return rsp, nil
	}
	return packPostListPageResp(req.GetPage(), hasMore, req.GetPageType(), category.CategoryName, req.GetName(), misc.MapToList(postIDs, postMetas)), nil

}
