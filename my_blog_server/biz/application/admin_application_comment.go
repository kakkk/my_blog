package application

import (
	"context"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/factory"
	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/infra/pkg/log"
)

func (a *AdminApplication) GetCommentList(ctx context.Context, req *api.GetCommentListAdminAPIRequest) (*api.GetCommentListAdminAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"page":  req.GetPage(),
		"limit": req.GetLimit(),
	})
	comments := factory.NewComments()
	// 先分页获取评论
	total, err := comments.GetByPage(ctx, req.Page, req.Limit)
	if err != nil {
		logger.Errorf("get comments by page err:%v", err)
		return nil, err
	}
	// 获取回复的评论的内容
	err = comments.GetParentsContent(ctx)
	if err != nil {
		logger.Errorf("get comments parent content err:%v", err)
		return nil, err
	}
	// 拿所有postID
	postIDs := comments.GetPostIDs()
	// 获取文章
	posts := factory.GetContentFactory().NewArticles()
	err = posts.MGetMyIDs(ctx, postIDs)
	if err != nil {
		logger.Errorf("mget posts by id err:%v", err)
	}
	// 拿标题
	postID2Title := posts.GetID2Title()
	// 设置评论的标题
	comments.SetArticleTitle(postID2Title)
	// 返回
	commentList := comments.ToDTOs().ToResponseAdminCommentList()
	return &api.GetCommentListAdminAPIResponse{
		Pagination: &api.Pagination{
			Page:    req.GetPage(),
			Limit:   req.GetLimit(),
			HasMore: true,
			Total:   &total,
		},
		Comments: commentList,
		BaseResp: api.NewBaseResp(),
	}, nil
}

func (a *AdminApplication) DeleteComment(ctx context.Context, req *api.DeleteCommentAdminAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"id": req.GetCommentID(),
	})
	comment := factory.NewCommentByID(req.GetCommentID())
	err := comment.Delete(ctx)
	if err != nil {
		logger.Errorf("delete comment err:%v", err)
		return nil, err
	}
	logger.Infof("delete comment success")
	return &api.CommonResponse{
		BaseResp: api.NewBaseResp(),
	}, nil
}
