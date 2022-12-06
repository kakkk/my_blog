package facade

import (
	"context"
	"net/http"

	"my_blog/biz/mock"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func PostPage(ctx context.Context, c *app.RequestContext) (resp *page.PostPageResp, code int) {
	id := c.Param("id")
	return mock.PostPageMocker(id), http.StatusOK
}
