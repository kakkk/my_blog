package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/api"
)

func CreateTagAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreateTagAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().CreateTag(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)

}

func UpdateTagAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdateTagAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().UpdateTag(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func DeleteTagAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeleteTagAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().DeleteTag(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetTagListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetTagListAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().GetTagList(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
