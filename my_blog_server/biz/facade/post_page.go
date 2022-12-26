package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func PostPage(ctx context.Context, c *app.RequestContext) (int, string, *page.PostPageResponse) {
	tmpl := "index.tmpl"
	req := &page.PostPageRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		rsp := &page.PostPageResponse{
			Meta: resp.NewNotFoundErrorMeta(),
		}
		return http.StatusNotFound, tmpl, rsp
	}
	code, rsp := service.PostPage(ctx, req)
	return code, tmpl, rsp
}
