package facade

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/consts"
	"my_blog/biz/infra/pkg/resp"
)

func TagsPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	rsp, pErr := application.GetBlogApplication().TagsPage(ctx)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func CategoriesPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	rsp, pErr := application.GetBlogApplication().CategoriesPage(ctx)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
