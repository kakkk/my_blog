package facade

import (
	"context"
	"net/http"

	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func SearchPage(ctx context.Context, c *app.RequestContext) (resp *page.BasicPageResp, code int) {
	return mock.SearchMocker(), http.StatusOK
}
