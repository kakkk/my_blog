package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func LoginAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.LoginRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := service.LoginAPI(ctx, req)
	if err != nil {
		return http.StatusOK, resp.NewInternalErrorResp()
	}

	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetUserInfoAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	rsp, err := service.GetUserInfoAPI(ctx)
	if err != nil {
		return http.StatusOK, resp.NewInternalErrorResp()
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
