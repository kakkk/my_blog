package facade

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"my_blog/biz/application"
	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/resp"
)

func GetCommentListAPI(ctx context.Context, c *app.RequestContext) (int, *resp.APIResponse) {
	var err error
	var req api.GetCommentListAPIRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		return http.StatusBadRequest, resp.NewParameterErrorResp()
	}

	rsp, err := application.GetBlogApplication().GetCommentList(ctx, &req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
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

	rsp, err := application.GetBlogApplication().CommentArticle(ctx, &req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}

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

	rsp, err := application.GetBlogApplication().ReplyComment(ctx, &req)
	if err != nil {
		return resp.NewErrorAPIResponse(err)
	}
	return http.StatusOK, resp.NewAPIResponse(rsp)
}
