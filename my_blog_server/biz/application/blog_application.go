package application

import (
	"context"
	"fmt"

	"my_blog/biz/consts"
	"my_blog/biz/domain/service"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/errorx"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/page"
)

type BlogApplication struct{}

func GetBlogApplication() *BlogApplication {
	return &BlogApplication{}
}

func (app *BlogApplication) GetArchives(ctx context.Context) (rsp *page.ArchivesPageResp, pErr *errorx.PageError) {
	rsp = &page.ArchivesPageResp{
		Meta: &page.PageMeta{
			Title:       fmt.Sprintf("文章归档 - %v", config.GetBlogName()),
			Description: "文章归档",
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			PageType:    page.PageTypeArchives,
		},
	}
	archives := service.GetIndexService().GetArchives(ctx)
	rsp.PostArchives = archives
	return rsp, nil
}

func (app *BlogApplication) Search(ctx context.Context, req *api.SearchAPIRequest) (rsp *api.SearchAPIResponse, err error) {
	logger := log.GetLoggerWithCtx(ctx).WithField("query", req.GetQuery())
	rsp = api.NewSearchAPIResponse()
	defer misc.Recover(ctx, func() {
		err = consts.ErrInternalServerError
	})()

	result, err := service.GetIndexService().SearchArticleByTitle(req.GetQuery(), 10)
	if err != nil {
		logger.Errorf("query error:[%v]", err)
		return nil, fmt.Errorf("query error:[%v]", err)
	}
	var res []*api.SearchResultItem

	for _, data := range result {
		res = append(res, &api.SearchResultItem{
			Link:     fmt.Sprintf("%v/archives/%v", config.GetSiteConfig().SiteDomain, data.ID),
			Title:    data.Title,
			Abstract: "",
		})
	}
	rsp.BaseResp = resp.NewSuccessBaseResp()
	rsp.Results = res
	return rsp, nil
}
