package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/errorx"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/service"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/hertz/pkg/app"
)

func IndexPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	req := &page.PostListPageRequest{
		Page:     thrift.Int64Ptr(1),
		PageType: thrift.StringPtr(page.PageTypeIndex),
	}
	rsp, pErr := service.PostListByPage(ctx, req)
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
	rsp, pErr := service.PostListByPage(ctx, req)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func CategoryPostPage(ctx context.Context, c *app.RequestContext) (resp *page.PostListPageResp, code int) {
	return mock.PostListPageRespMocker(page.PageTypeCategoryPostList, "测试分类", "", "2", "test_category"), http.StatusOK
}

func CategoryPostByPaginationPage(ctx context.Context, c *app.RequestContext) (resp *page.PostListPageResp, code int) {
	p := c.Param("page")
	if p == "1" {
		return CategoryPostPage(ctx, c)
	}
	return mock.PostListPageRespMocker(page.PageTypeCategoryPostList, "测试分类", "1", "3", "test_category"), http.StatusOK
}

func TagPostPage(ctx context.Context, c *app.RequestContext) (resp *page.PostListPageResp, code int) {
	name := c.Param("name")
	return mock.PostListPageRespMocker(page.PageTypeTagPostList, name, "", "2", ""), http.StatusOK
}

func TagPostByPaginationPage(ctx context.Context, c *app.RequestContext) (resp *page.PostListPageResp, code int) {
	p := c.Param("page")
	if p == "1" {
		return TagPostPage(ctx, c)
	}
	name := c.Param("name")
	return mock.PostListPageRespMocker(page.PageTypeTagPostList, name, "1", "3", ""), http.StatusOK
}
