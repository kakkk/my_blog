package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/api"
)

func LoginAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.LoginRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().Login(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetUserInfoAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	rsp, err := application.GetAdminApplication().GetUserInfo(ctx)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
