package resp

import (
	"errors"
	"fmt"

	"my_blog/biz/common/config"
	"my_blog/biz/common/errorx"
	"my_blog/biz/model/blog/page"
)

type IPageResponse interface {
	GetMeta() (v *page.PageMeta)
}

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

func NewSomethingWrongErrorPageResp() *page.BasicPageResp {
	return &page.BasicPageResp{
		Meta: NewSomethingWrongErrorMeta(),
	}
}

func NewInternalErrorMeta() *page.PageMeta {
	return &page.PageMeta{
		Title:       fmt.Sprintf("Internal Error - %v", config.GetBlogName()),
		Description: "Internal Error",
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    page.PageTypeError,
		ErrorCode:   "500",
	}
}

func NewNotFoundErrorMeta() *page.PageMeta {
	return &page.PageMeta{
		Title:       fmt.Sprintf("Not Found - %v", config.GetBlogName()),
		Description: "Not Found",
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    page.PageTypeError,
		ErrorCode:   "404",
	}
}

func NewSomethingWrongErrorMeta() *page.PageMeta {
	return &page.PageMeta{
		Title:       fmt.Sprintf("Something Wrong - %v", config.GetBlogName()),
		Description: "Something Wrong",
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    page.PageTypeError,
		ErrorCode:   "( •︠ˍ•︡ )",
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

func NewBasePageMeta(pageType string) *page.PageMeta {
	return &page.PageMeta{
		Title:       "",
		Description: config.GetBlogDescription(),
		CDNDomain:   config.GetSiteConfig().CDNDomain,
		SiteDomain:  config.GetSiteConfig().SiteDomain,
		PageType:    pageType,
	}
}

func PackPageResponse(rsp IPageResponse, pErr *errorx.PageError, tmpl string) (int, string, IPageResponse) {
	if pErr.IsError() {
		if errors.Is(pErr, errorx.PageErrInternalError) {
			return pErr.GetStatusCode(), tmpl, NewInternalErrorPageResp()
		}
		if errors.Is(pErr, errorx.PageErrNotFound) {
			return pErr.GetStatusCode(), tmpl, NewNotFoundErrorPageResp()
		}
		if errors.Is(pErr, errorx.PageErrFail) {
			return pErr.GetStatusCode(), tmpl, NewSomethingWrongErrorPageResp()
		}
		return pErr.GetStatusCode(), tmpl, NewInternalErrorPageResp()
	}

	return pErr.GetStatusCode(), tmpl, rsp
}
