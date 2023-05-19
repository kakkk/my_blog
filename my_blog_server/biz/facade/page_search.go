package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/config"
	"my_blog/biz/model/blog/page"

	"github.com/cloudwego/hertz/pkg/app"
)

func SearchPage(ctx context.Context, c *app.RequestContext) (resp *page.BasicPageResp, code int) {
	return &page.BasicPageResp{
		Meta: &page.PageMeta{
			Title:       "Search",
			Description: "search",
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeSearch,
		},
	}, http.StatusOK
}
