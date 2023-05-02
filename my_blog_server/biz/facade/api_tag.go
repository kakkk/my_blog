package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func CreateTagAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreateTagAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := service.CreateTagAPI(ctx, req)
	if err != nil {
		return http.StatusOK, resp.NewInternalErrorResp()
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)

}

func UpdateTagAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdateTagAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := service.UpdateTagAPI(ctx, req)
	if err != nil {
		return http.StatusOK, resp.NewInternalErrorResp()
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func DeleteTagAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeleteTagAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := service.DeleteTagAPI(ctx, req)
	if err != nil {
		return http.StatusOK, resp.NewInternalErrorResp()
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetTagListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetTagListAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := service.GetTagListAPI(ctx, req)
	if err != nil {
		return http.StatusOK, resp.NewInternalErrorResp()
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
