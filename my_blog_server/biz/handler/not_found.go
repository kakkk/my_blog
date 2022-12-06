package handler

import (
	"context"
	"net/http"

	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func NotFoundHandler(ctx context.Context, c *app.RequestContext) {
	c.HTML(http.StatusNotFound, "index.tmpl", page.BasicPageResp{
		Meta: &page.PageMeta{
			Title:       "404",
			Description: "not found",
			SiteDomain:  "http://127.0.0.1:8888",
			PageType:    page.PageTypeError,
			ErrorCode:   "404",
		},
	})
}
