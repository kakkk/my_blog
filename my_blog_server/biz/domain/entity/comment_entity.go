package entity

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/idgen"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/common"
)

type Comment struct {
	ID       int64
	PostID   int64
	ReplyID  int64
	ParentID int64
	Nickname string
	Email    string
	Website  string
	Content  string
	Status   common.CommentStatus
	CreateAt time.Time
}

func (c *Comment) Create(ctx context.Context) error {
	c.ID = idgen.GenID()
	c.Status = c.Review(ctx)
	err := repo.GetCommentRepo().Create(mysql.GetDB(ctx), c.GetDTO())
	if err != nil {
		return fmt.Errorf("repo create fail: %w", err)
	}
	return nil
}

func (c *Comment) ReplyTo(ctx context.Context, replyID int64) error {
	// 检查comment是否存在
	replyComment, err := repo.GetCommentRepo().Cache().Get(ctx, replyID)
	if err != nil {
		// 不存在会返回not found
		return err
	}

	// 若回复的评论是子评论，父评论为该评论的父评论
	if replyComment.ParentID != 0 {
		c.ParentID = replyComment.ParentID
	} else {
		// 若回复的评论是父评论，父评论为该评论
		c.ParentID = replyComment.ID
	}

	c.ID = idgen.GenID()
	err = repo.GetCommentRepo().Create(mysql.GetDB(ctx), c.GetDTO())
	if err != nil {
		return fmt.Errorf("repo create fail: %w", err)
	}
	return nil
}

// Review 评论审核，当前直接通过
func (c *Comment) Review(ctx context.Context) common.CommentStatus {
	// TODO：评论审核相关逻辑
	return common.CommentStatus_Approved
}

func (c *Comment) GetDTO() *dto.Comment {
	return &dto.Comment{
		ID:       c.ID,
		PostID:   c.PostID,
		ReplyID:  c.ReplyID,
		ParentID: c.ParentID,
		Nickname: c.Nickname,
		Email:    c.Email,
		Website:  c.Website,
		Content:  c.Content,
		CreateAt: c.CreateAt,
	}
}

func NewCommentByDTO(c *dto.Comment) *Comment {
	return &Comment{
		ID:       c.ID,
		PostID:   c.PostID,
		ReplyID:  c.ReplyID,
		ParentID: c.ParentID,
		Nickname: c.Nickname,
		Email:    c.Email,
		Website:  c.Website,
		Content:  c.Content,
		CreateAt: c.CreateAt,
	}
}
