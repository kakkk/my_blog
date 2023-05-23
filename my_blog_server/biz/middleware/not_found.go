package middleware

import (
	"context"
	"net/http"

	"my_blog/biz/common/config"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func NotFoundMW() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
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
}
