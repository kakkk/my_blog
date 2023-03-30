package facade

import (
	"context"
	"net/http"

	"my_blog/biz/common/resp"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
)

func CreatePostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreatePostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.CreatePostAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetPostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetPostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.GetPostAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdatePostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdatePostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.UpdatePostAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdatePostStatusAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdatePostStatusAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.UpdatePostStatusAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetPostListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetPostListAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.GetPostListAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func DeletePostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeletePostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.DeletePostAPI(ctx, req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
