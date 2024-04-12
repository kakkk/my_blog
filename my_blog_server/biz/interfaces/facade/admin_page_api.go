package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/infra/pkg/resp"
)

func CreatePageAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreatePageAPIRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().CreatePage(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetPageAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetPageAPIRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().GetPageByID(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdatePageAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdatePageAPIRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().UpdatePage(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetPageListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	rsp, err := application.GetAdminApplication().GetPageList(ctx)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func DeletePageAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeletePageAPIRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().DeletePage(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
