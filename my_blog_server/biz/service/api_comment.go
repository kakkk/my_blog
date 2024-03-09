package service

import (
	"context"
	"errors"
	"sort"

	"github.com/sirupsen/logrus"

	"my_blog/biz/common/resp"
	"my_blog/biz/consts"
	"my_blog/biz/dto"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/repository/model"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/storage"
)

func CommentArticleAPI(ctx context.Context, req *api.CommentArticleAPIRequest) (rsp *api.CommentArticleAPIResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"article_id": req.GetArticleID(),
		"nickname":   req.GetNickname(),
	})
	rsp = &api.CommentArticleAPIResponse{}
	defer misc.Recover(ctx, func() {
		logger.Errorf("recovered")
		rsp.BaseResp = resp.NewInternalErrorBaseResp()
		return
	})()
	// 检查article是否存在
	_, err := storage.GetArticleMetaStorage().Get(ctx, req.GetArticleID())
	if errors.Is(err, consts.ErrRecordNotFound) {
		logger.Warnf("article not found")
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
	}
	comment := &model.Comment{
		PostID:   req.GetArticleID(),
		Nickname: req.GetNickname(),
		Email:    req.GetEmail(),
		Website:  req.GetWebsite(),
		Content:  req.GetContent(),
		Status:   common.CommentStatus_Approved,
	}
	comment, err = mysql.CreateComment(mysql2.GetDB(ctx), comment)
	if err != nil {
		logger.Error("create comment fail")
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
	}

	rsp = &api.CommentArticleAPIResponse{
		ID:            comment.ID,
		CommentStatus: common.CommentStatus_Approved,
		BaseResp:      resp.NewSuccessBaseResp(),
	}
	return rsp
}

func ReplyCommentAPI(ctx context.Context, req *api.ReplyCommentAPIRequest) (rsp *api.ReplyCommentAPIResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"article_id": req.GetArticleID(),
		"nickname":   req.GetNickname(),
		"comment_id": req.GetReplyID(),
	})
	rsp = &api.ReplyCommentAPIResponse{}
	defer misc.Recover(ctx, func() {
		logger.Errorf("recovered")
		rsp.BaseResp = resp.NewInternalErrorBaseResp()
		return
	})()

	// 检查article是否存在
	_, err := storage.GetArticleMetaStorage().Get(ctx, req.GetArticleID())
	if errors.Is(err, consts.ErrRecordNotFound) {
		logger.Warnf("article not found")
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
	}

	// 检查comment是否存在
	replyComment, err := storage.GetCommentStorageStorage().Get(ctx, req.GetReplyID())
	if errors.Is(err, consts.ErrRecordNotFound) {
		logger.Warnf("comment not found")
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
	}

	// 父评论ID
	var parentID int64
	// 若回复的评论是子评论，父评论为该评论的父评论
	if replyComment.ParentID != 0 {
		parentID = replyComment.ParentID
	} else {
		// 若回复的评论是父评论，父评论为该评论
		parentID = replyComment.ID
	}

	comment := &model.Comment{
		PostID:   req.GetArticleID(),
		Nickname: req.GetNickname(),
		Email:    req.GetEmail(),
		Website:  req.GetWebsite(),
		Content:  req.GetContent(),
		ReplyID:  req.GetReplyID(),
		ParentID: parentID,
		Status:   common.CommentStatus_Approved,
	}
	comment, err = mysql.CreateComment(mysql2.GetDB(ctx), comment)
	if err != nil {
		logger.Error("create comment fail")
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
	}

	rsp = &api.ReplyCommentAPIResponse{
		ID:            comment.ID,
		CommentStatus: common.CommentStatus_Approved,
		BaseResp:      resp.NewSuccessBaseResp(),
	}
	return rsp
}

func GetCommentListAPI(ctx context.Context, req *api.GetCommentListAPIRequest) (rsp *api.GetCommentListAPIResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"article_id": req.GetArticleID(),
	})
	rsp = api.NewGetCommentListAPIResponse()
	// 检查article是否存在
	_, err := storage.GetArticleMetaStorage().Get(ctx, req.GetArticleID())
	if errors.Is(err, consts.ErrRecordNotFound) {
		logger.Warnf("article not found")
		rsp.BaseResp = resp.NewFailBaseResp()
		return rsp
	}

	ids := storage.GetPostCommentIDsStorage().Get(ctx, req.GetArticleID())
	commentMap := storage.GetCommentStorageStorage().MGet(ctx, ids)
	list := packCommentList(commentMap)
	rsp = &api.GetCommentListAPIResponse{
		Comments: list,
		HasMore:  false,
		BaseResp: resp.NewSuccessBaseResp(),
	}
	return rsp
}

func packCommentList(comments map[int64]*dto.Comment) []*api.CommentListItem {
	if len(comments) == 0 {
		return make([]*api.CommentListItem, 0)
	}
	var commentList []*api.CommentListItem
	pIDToComments := make(map[int64][]*dto.Comment)
	// 第一次遍历，找到所有父评论和子评论
	for _, comment := range comments {
		if comment.ParentID == 0 {
			item := comment.ToAPIModel()
			commentList = append(commentList, &api.CommentListItem{Comment: item})
			continue
		}
		_, ok := pIDToComments[comment.ParentID]
		if !ok {
			pIDToComments[comment.ParentID] = []*dto.Comment{comment}
			continue
		}
		pIDToComments[comment.ParentID] = append(pIDToComments[comment.ParentID], comment)
	}
	// 子评论排序
	for pID, c := range pIDToComments {
		sort.Slice(c, func(i, j int) bool {
			return c[i].ID > c[j].ID
		})
		pIDToComments[pID] = c
	}
	// 父评论排序
	sort.Slice(commentList, func(i, j int) bool {
		return commentList[i].Comment.ID > commentList[j].Comment.ID
	})
	// 组装
	for _, item := range commentList {
		list := pIDToComments[item.Comment.ID]
		replies := make([]*api.Comment, 0, len(list))
		for _, comment := range list {
			replyName := comments[comment.ReplyID].Nickname
			replies = append(replies, comment.ToAPIModelWithAtUser(replyName))
		}
		item.Replies = replies
	}
	return commentList
}
