package facade

import (
	"context"
	"net/http"

	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func TagsPage(ctx context.Context, c *app.RequestContext) (resp *page.TermsPageResp, code int) {
	return mock.TagsMocker(), http.StatusOK
}

func CategoriesPage(ctx context.Context, c *app.RequestContext) (resp *page.TermsPageResp, code int) {
	return mock.CategoriesMocker(), http.StatusOK
}
