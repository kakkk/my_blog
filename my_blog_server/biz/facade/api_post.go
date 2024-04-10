package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/model/blog/api"
)

func CreatePostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.CreatePostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().CreatePost(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetPostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetPostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().GetPostByID(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdatePostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdatePostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().UpdatePost(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func UpdatePostStatusAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.UpdatePostStatusAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().UpdatePostStatus(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func GetPostListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.GetPostListAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().GetPostList(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func DeletePostAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	req := &api.DeletePostAPIRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().DeletePost(ctx, req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
