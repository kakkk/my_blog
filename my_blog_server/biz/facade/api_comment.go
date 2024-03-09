package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/common/resp"
	"my_blog/biz/infra/misc"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/service"
)

// ========管理员接口=========

func GetCommentListAdminAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.GetCommentListAdminAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	var rsp api.GetCommentListAdminAPIResponse

	//rsp := service.GetCommentListAdminAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(&rsp)
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
	var rsp api.CommonResponse

	//rsp := service.DeleteCommentAdminAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(&rsp)
}

// ========用户侧接口=========

func GetCommentListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.GetCommentListAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.GetCommentListAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func CommentArticleAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.CommentArticleAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	// 参数校验
	if req.GetNickname() == "" ||
		req.GetEmail() == "" ||
		req.GetContent() == "" ||
		req.GetArticleID() == 0 {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	if !misc.CheckIsEmail(req.GetEmail()) || !misc.CheckIsURL(req.GetWebsite()) {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.CommentArticleAPI(ctx, &req)

	return http.StatusOK, resp.NewAPIResponse(rsp)
}

func ReplyCommentAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.ReplyCommentAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	// 参数校验
	if req.GetNickname() == "" ||
		req.GetEmail() == "" ||
		req.GetContent() == "" ||
		req.GetArticleID() == 0 ||
		req.GetReplyID() == 0 {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}
	if !misc.CheckIsEmail(req.GetEmail()) || !misc.CheckIsURL(req.GetWebsite()) {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp := service.ReplyCommentAPI(ctx, &req)
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
