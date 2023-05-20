package service

import (
	"context"
	"fmt"

	"my_blog/biz/common/config"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/repository/index"
)

func SearchAPI(ctx context.Context, req *api.SearchAPIRequest) (rsp *api.SearchAPIResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithField("query", req.GetQuery())
	rsp = api.NewSearchAPIResponse()
	defer utils.Recover(ctx, func() {
		logger.Errorf("recovered")
		rsp.BaseResp = resp.NewInternalErrorBaseResp()
		return
	})()

	result, err := index.GetArticleIndex().SearchByTitle(req.GetQuery(), 10)
	if err != nil {
		logger.Errorf("query error:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
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
	return rsp
}
