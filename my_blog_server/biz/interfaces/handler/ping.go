// Code generated by hertz generator.

package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/hertz_gen/blog/page"
	"my_blog/biz/infra/config"
)

// Ping .
func Ping(ctx context.Context, c *app.RequestContext) {
	c.Response.Header.Set("x-rsp", "pong")
	c.HTML(http.StatusNotFound, "index.tmpl", page.BasicPageResp{
		Meta: &page.PageMeta{
			Title:       "404",
			Description: "404 not found",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			PageType:    page.PageTypeError,
			ErrorCode:   "404",
		},
	})
}
