package facade

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/resp"
	"my_blog/biz/service"
)

func TagsPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	rsp, pErr := service.TagsPage(ctx)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}

func CategoriesPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	rsp, pErr := service.CategoriesPage(ctx)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
