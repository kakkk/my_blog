package application

import (
	"context"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/factory"
	"my_blog/biz/domain/repo"
	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
)

// CommentArticle 评论文章
func (app *BlogApplication) CommentArticle(ctx context.Context, req *api.CommentArticleAPIRequest) (*api.CommentArticleAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"article_id": req.GetArticleID(),
		"nickname":   req.GetNickname(),
	})
	// 检查article是否存在
	_, err := repo.GetContentRepo().Cache().GetArticleMeta(ctx, req.GetArticleID())
	if err != nil {
		logger.Warnf("article not found")
		// 不存在会返回notfound
		return nil, err
	}
	// 创建评论
	comment := factory.NewCommentByDTO(&dto.Comment{
		ArticleID: req.GetArticleID(),
		Nickname:  req.GetNickname(),
		Email:     req.GetEmail(),
		Website:   req.GetWebsite(),
		Content:   req.GetContent(),
	})
	err = comment.Create(ctx)
	if err != nil {
		logger.Errorf("create comment fail: %v", err)
		return nil, err
	}
	// 刷新缓存
	repo.GetCommentRepo().Cache().RefreshArticleComments(ctx, req.GetArticleID())
	// 获取评论列表
	comments := repo.GetCommentRepo().Cache().GetArticleComments(ctx, req.GetArticleID())
	logger.Infof("create comment success")
	return &api.CommentArticleAPIResponse{
		ID:            comment.ID,
		CommentStatus: comment.Status,
		Comments:      dto.Comments(comments).ToResponseCommentList(),
		BaseResp:      resp.NewSuccessBaseResp(),
	}, nil
}

// ReplyComment 回复文章
func (app *BlogApplication) ReplyComment(ctx context.Context, req *api.ReplyCommentAPIRequest) (*api.CommentArticleAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"article_id": req.GetArticleID(),
		"nickname":   req.GetNickname(),
		"reply_to":   req.GetReplyID(),
	})
	// 检查article是否存在
	_, err := repo.GetContentRepo().Cache().GetArticleMeta(ctx, req.GetArticleID())
	if err != nil {
		logger.Warnf("article not found")
		// 不存在会返回notfound
		return nil, err
	}
	// 回复评论
	comment := factory.NewCommentByDTO(&dto.Comment{
		ArticleID: req.GetArticleID(),
		Nickname:  req.GetNickname(),
		Email:     req.GetEmail(),
		Website:   req.GetWebsite(),
		Content:   req.GetContent(),
	})
	err = comment.ReplyTo(ctx, req.GetReplyID())
	if err != nil {
		logger.Errorf("reply comment fail: %v", err)
		return nil, err
	}
	// 刷新缓存
	repo.GetCommentRepo().Cache().RefreshArticleComments(ctx, req.GetArticleID())
	// 获取评论列表
	comments := repo.GetCommentRepo().Cache().GetArticleComments(ctx, req.GetArticleID())
	logger.Infof("reply comment success")
	return &api.CommentArticleAPIResponse{
		ID:            comment.ID,
		CommentStatus: comment.Status,
		Comments:      dto.Comments(comments).ToResponseCommentList(),
		BaseResp:      resp.NewSuccessBaseResp(),
	}, nil
}

// GetCommentList 获取文章评论列表
func (app *BlogApplication) GetCommentList(ctx context.Context, req *api.GetCommentListAPIRequest) (*api.GetCommentListAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"article_id": req.GetArticleID(),
	})
	// 检查article是否存在
	_, err := repo.GetContentRepo().Cache().GetArticleMeta(ctx, req.GetArticleID())
	if err != nil {
		logger.Warnf("article not found")
		// 不存在会返回notfound
		return nil, err
	}
	// 获取评论
	comments := repo.GetCommentRepo().Cache().GetArticleComments(ctx, req.GetArticleID())
	return &api.GetCommentListAPIResponse{
		Comments: dto.Comments(comments).ToResponseCommentList(),
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil

}
