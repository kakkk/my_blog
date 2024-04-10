package facade

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/consts"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/page"
)

func IndexPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		Page:     thrift.Int64Ptr(1),
		PageType: thrift.StringPtr(page.PageTypeIndex),
	}
	rsp, pErr := application.GetBlogApplication().PostListPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func IndexByPaginationPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		PageType: thrift.StringPtr(page.PageTypePostList),
	}
	err := c.BindAndValidate(req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := application.GetBlogApplication().PostListPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func CategoryPostPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		Page:     thrift.Int64Ptr(1),
		PageType: thrift.StringPtr(page.PageTypeCategoryPostList),
	}
	err := c.BindAndValidate(req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := application.GetBlogApplication().CategoryPostListPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func CategoryPostByPaginationPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		PageType: thrift.StringPtr(page.PageTypeCategoryPostList),
	}
	err := c.BindAndValidate(req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := application.GetBlogApplication().CategoryPostListPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func TagPostPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		Page:     thrift.Int64Ptr(1),
		PageType: thrift.StringPtr(page.PageTypeTagPostList),
	}
	err := c.BindAndValidate(req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := application.GetBlogApplication().TagPostListPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func TagPostByPaginationPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		PageType: thrift.StringPtr(page.PageTypeTagPostList),
	}
	err := c.BindAndValidate(req)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Warnf("parameter error:[%v]", err)
		return resp.PackPageResponse(nil, errorx.NewNotFoundErrPageError(), consts.IndexTmpl)
	}
	rsp, pErr := application.GetBlogApplication().TagPostListPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
