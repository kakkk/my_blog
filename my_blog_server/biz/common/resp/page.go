package resp

import (
	"my_blog/biz/common/config"
	"my_blog/biz/model/blog/page"
)

func NewNotFoundErrorPageResp() *page.BasicPageResp {
	return &page.BasicPageResp{
		Meta: NewNotFoundErrorMeta(),
	}
}

func NewInternalErrorPageResp() *page.BasicPageResp {
	return &page.BasicPageResp{
		Meta: NewInternalErrorMeta(),
	}
}

func NewInternalErrorMeta() *page.PageMeta {
	return &page.PageMeta{
		Title:       "500 Internal Error",
		Description: "500 Internal Error",
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    page.PageTypeError,
		ErrorCode:   "500",
	}
}

func NewNotFoundErrorMeta() *page.PageMeta {
	return &page.PageMeta{
		Title:       "404 Not Found",
		Description: "404 Not Found",
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    page.PageTypeError,
		ErrorCode:   "404",
	}
}

func NewSuccessPageMeta(title string, description string, pageType string) *page.PageMeta {
	return &page.PageMeta{
		Title:       title,
		Description: description,
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    pageType,
	}
}
