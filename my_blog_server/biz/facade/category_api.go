package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func CreateCategoryAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreateCategoryAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.CreateCategoryAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp.GetBaseResp(), rsp)
}

func UpdateCategoryAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdateCategoryAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.UpdateCategoryAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp.GetBaseResp(), rsp)
}

func DeleteCategoryAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeleteCategoryAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.DeleteCategoryAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp.GetBaseResp(), rsp)
}

func UpdateCategoryOrderAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdateCategoryOrderAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.UpdateCategoryOrderAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp.GetBaseResp(), rsp)
}

func GetCategoryListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	rsp := service.GetCategoryListAPI(ctx)
	return http.StatusOK, resp.NewAPIResponse(rsp.GetBaseResp(), rsp)
}
