package facade

import (
	"context"
	"my_blog/biz/common/consts"
	"my_blog/biz/common/resp"
	"my_blog/biz/service"
	"net/http"

	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func TagsPage(ctx context.Context, c *app.RequestContext) (resp *page.TermsPageResp, code int) {
	return mock.TagsMocker(), http.StatusOK
}

func CategoriesPage(ctx context.Context, c *app.RequestContext) (int, string, resp.IPageResponse) {
	rsp, pErr := service.CategoriesPage(ctx)
	return resp.PackPageResponse(rsp, pErr, consts.IndexTmpl)
}
