package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/api"
)

func CreateCategoryAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreateCategoryAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().CreateCategory(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdateCategoryAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdateCategoryAPIRequest{}
	err := c.BindAndValidate(&req)

	rsp, err := application.GetAdminApplication().UpdateCategory(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func DeleteCategoryAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeleteCategoryAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().DeleteCategory(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdateCategoryOrderAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdateCategoryOrderAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	// TODO
	rsp := &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetCategoryListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	rsp, err := application.GetAdminApplication().GetCategoryList(ctx)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
