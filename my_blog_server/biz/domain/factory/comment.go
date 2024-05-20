package factory

import (
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
)

func NewCommentByDTO(c *dto.Comment) *entity.Comment {
	return &entity.Comment{
		ID:       c.ID,
		PostID:   c.ArticleID,
		ReplyID:  c.ReplyID,
		ParentID: c.ParentID,
		Nickname: c.Nickname,
		Email:    c.Email,
		Website:  c.Website,
		Content:  c.Content,
		CreateAt: c.CreateAt,
	}
}

func NewComments() *entity.Comments {
	return &entity.Comments{
		Comments: make([]*entity.Comment, 0),
	}
}
