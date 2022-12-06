package facade

import (
	"context"
	"net/http"

	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func IndexPage(ctx context.Context, c *app.RequestContext) (resp *page.PostListPageResp, code int) {
	return mock.PostListPageRespMocker(page.PageTypeIndex, "", "", "2", ""), http.StatusOK
}

func IndexByPaginationPage(ctx context.Context, c *app.RequestContext) (resp *page.PostListPageResp, code int) {
	p := c.Param("page")
	if p == "1" {
		return IndexPage(ctx, c)
	}
	return mock.PostListPageRespMocker(page.PageTypePostList, "", "1", "3", ""), http.StatusOK
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
