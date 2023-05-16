package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/common/errorx"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/storage"
)

func TagPostListByPage(ctx context.Context, req *page.PostListPageRequest) (rsp *page.PostListPageResp, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"name": req.GetName(),
		"page": req.GetPage(),
	})
	rsp = &page.PostListPageResp{
		Meta:     resp.NewBasePageMeta(page.PageTypeTagPostList),
		PostList: []*page.PostItem{},
	}
	defer utils.Recover(ctx, func() {
		pErr = errorx.NewInternalErrPageError()
		return
	})()
	// 获取标签id
	categoryID, err := storage.GetTagNameIDStorage().Get(ctx, req.GetName())
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("tag not found")
			return nil, errorx.NewNotFoundErrPageError()
		}
		logger.Errorf("get tag_id error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	rsp.Meta.Title = fmt.Sprintf("标签 %v 下的文章 - %v", req.GetName(), config.GetBlogName())
	rsp.Name = req.GetName()
	rsp.Slug = req.GetName()
	// 获取标签下所有文章id
	postIDs, err := storage.GetTagPostListStorage().Get(ctx, categoryID)
	if err != nil {
		if !errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("no article")
			return rsp, nil
		}
		logger.Errorf("get tag_post_list error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	postIDs, hasMore := getIDsByPage(postIDs, req.GetPage())
	if len(postIDs) == 0 {
		logger.Warnf("no article")
		return rsp, nil
	}
	postMetas := storage.GetPostMetaStorage().MGet(ctx, postIDs)
	if len(postMetas) == 0 {
		logger.Warnf("no article")
		return rsp, nil
	}
	return packPostListPageResp(req.GetPage(), hasMore, req.GetPageType(), req.GetName(), req.GetName(), utils.MapToList(postIDs, postMetas)), nil

}
