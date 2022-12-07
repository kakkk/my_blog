package resp

import (
	"my_blog/biz/common/config"
	"my_blog/biz/model/blog/page"
)

func GetInternalErrorPageResp() *page.BasicPageResp {
	return &page.BasicPageResp{
		Meta: GetInternalErrorMeta(),
	}
}

func GetInternalErrorMeta() *page.PageMeta {
	return &page.PageMeta{
		Title:       "Internal Error",
		Description: "Internal Error",
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    "PAGE_ERROR",
		ErrorCode:   "500",
	}
}
