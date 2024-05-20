package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/infra/pkg/resp"
)

func GetCommentListAdminAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.GetCommentListAdminAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().GetCommentList(ctx, &req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func ReplyCommentAdminAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.ReplyCommentAdminAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	var rsp api.CommonResponse

	//rsp := service.ReplyCommentAdminAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(&rsp)
}

func UpdateCommentStatusAdminAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.UpdateCommentStatusAdminAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	var rsp api.CommonResponse

	//rsp := service.UpdateCommentStatusAdminAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(&rsp)
}

func DeleteCommentAdminAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.DeleteCommentAdminAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetAdminApplication().DeleteComment(ctx, &req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
