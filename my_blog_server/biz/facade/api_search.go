package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/service"
)

func SearchAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.SearchAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.SearchAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
