package service

import (
	"context"
	"fmt"

	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/storage"
)

func CategoriesPage(ctx context.Context) (rsp *page.TermsPageResp, pErr *errorx.PageError) {
	defer misc.Recover(ctx, func() {
		pErr = errorx.NewInternalErrPageError()
		return
	})()
	categoryList, err := storage.GetCategoryListStorage().Get(ctx)
	if err != nil {
		log.GetLoggerWithCtx(ctx).Errorf("get category list error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	title := fmt.Sprintf("分类 - %v - %v", config.GetBlogName(), config.GetBlogSubTitle())
	return &page.TermsPageResp{
		List: categoryList.ToPageCategoryListModel(),
		Meta: &page.PageMeta{
			Title:       title,
			Description: config.GetBlogDescription(),
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    page.PageTypeCategoryList,
		},
	}, nil
}
