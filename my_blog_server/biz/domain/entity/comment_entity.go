package entity

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/idgen"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/repository/mysql"
)

type Comment struct {
	ID            int64
	PostID        int64
	ReplyID       int64
	ParentID      int64
	Nickname      string
	Email         string
	Website       string
	Content       string
	ParentContent string
	ArticleTitle  string
	Status        common.CommentStatus
	CreateAt      time.Time
}

func (c *Comment) Create(ctx context.Context) error {
	c.ID = idgen.GenID()
	c.Status = c.Review(ctx)
	err := repo.GetCommentRepo().Create(mysql.GetDB(ctx), c.ToDTO())
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
	c.ReplyID = replyID

	c.ID = idgen.GenID()
	c.Status = c.Review(ctx)
	err = repo.GetCommentRepo().Create(mysql.GetDB(ctx), c.ToDTO())
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

func (c *Comment) fillByDTO(comment *dto.Comment) {
	c.ID = comment.ID
	c.PostID = comment.ArticleID
	c.ReplyID = comment.ReplyID
	c.ParentID = comment.ParentID
	c.ParentContent = comment.ParentContent
	c.ArticleTitle = comment.ArticleTitle
	c.Nickname = comment.Nickname
	c.Email = comment.Email
	c.Website = comment.Website
	c.Content = comment.Content
	c.CreateAt = comment.CreateAt
	c.Status = comment.Status
}

func (c *Comment) ToDTO() *dto.Comment {
	return &dto.Comment{
		ID:            c.ID,
		ArticleID:     c.PostID,
		ReplyID:       c.ReplyID,
		ParentID:      c.ParentID,
		ParentContent: c.ParentContent,
		ArticleTitle:  c.ArticleTitle,
		Nickname:      c.Nickname,
		Email:         c.Email,
		Website:       c.Website,
		Content:       c.Content,
		CreateAt:      c.CreateAt,
		Status:        c.Status,
	}
}

func (c *Comment) Get(ctx context.Context) error {
	comment, err := repo.GetCommentRepo().GetCommentByID(mysql.GetDB(ctx), c.ID)
	if err != nil {
		return err
	}
	c.fillByDTO(comment)
	return nil
}

func (c *Comment) Delete(ctx context.Context) error {
	err := c.Get(ctx)
	if err != nil {
		return fmt.Errorf("get comment fail: %w", err)
	}
	err = repo.GetCommentRepo().DeleteByID(mysql.GetDB(ctx), c.ID)
	if err != nil {
		return fmt.Errorf("delete comment fail: %w", err)
	}
	// 更新缓存
	repo.GetCommentRepo().Cache().RefreshArticleComments(ctx, c.PostID)
	return nil
}

type Comments struct {
	Comments []*Comment
}

func (c *Comments) GetByPage(ctx context.Context, page *int32, size *int32) (int64, error) {
	comments, err := repo.GetCommentRepo().GetListByPage(mysql.GetDB(ctx), page, size)
	if err != nil {
		return 0, fmt.Errorf("get comments fail: %w", err)
	}
	total, err := repo.GetCommentRepo().GetCount(mysql.GetDB(ctx))
	if err != nil {
		return 0, fmt.Errorf("get total comments fail: %w", err)
	}
	for _, comment := range comments {
		c.Comments = append(c.Comments, &Comment{
			ID:       comment.ID,
			PostID:   comment.ArticleID,
			ReplyID:  comment.ReplyID,
			ParentID: comment.ParentID,
			Nickname: comment.Nickname,
			Email:    comment.Email,
			Website:  comment.Website,
			Content:  comment.Content,
			CreateAt: comment.CreateAt,
		})
	}
	return total, nil
}

func (c *Comments) GetParentsContent(ctx context.Context) error {
	parentIDSet := make(map[int64]struct{})
	for _, comment := range c.Comments {
		if comment.ParentID != 0 {
			parentIDSet[comment.ParentID] = struct{}{}
		}
	}
	parents, err := repo.GetCommentRepo().MGetCommentsByID(mysql.GetDB(ctx), misc.MapKeys(parentIDSet))
	if err != nil {
		return err
	}
	for _, comment := range c.Comments {
		parent, ok := parents[comment.ParentID]
		if !ok {
			continue
		}
		comment.ParentContent = parent.Content
	}
	return nil
}

func (c *Comments) GetPostIDs() []int64 {
	postIDSet := make(map[int64]struct{})
	for _, comment := range c.Comments {
		postIDSet[comment.PostID] = struct{}{}
	}
	return misc.MapKeys(postIDSet)
}

func (c *Comments) SetArticleTitle(titles map[int64]string) {
	for _, comment := range c.Comments {
		title, ok := titles[comment.PostID]
		if ok {
			comment.ArticleTitle = title
			continue
		}
		title = "-"
	}
}

func (c *Comments) ToDTOs() dto.Comments {
	result := make([]*dto.Comment, 0, len(c.Comments))
	for _, comment := range c.Comments {
		result = append(result, comment.ToDTO())
	}
	return result
}
