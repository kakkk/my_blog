package dto

import (
	"time"

	"my_blog/biz/common/config"
	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/api"
)

type Comment struct {
	ID       int64     `json:"id"`
	PostID   int64     `json:"post_id"`
	ReplyID  int64     `json:"reply_id"`
	ParentID int64     `json:"parent_id"`
	Nickname string    `json:"nickname"`
	EmailMD5 string    `json:"email_md5"`
	Website  string    `json:"website"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
}

func NewCommentWithEntity(comment *entity.Comment) *Comment {
	return &Comment{
		ID:       comment.ID,
		PostID:   comment.PostID,
		ReplyID:  comment.ReplyID,
		ParentID: comment.ParentID,
		Nickname: comment.Nickname,
		EmailMD5: utils.SumStrMD5(comment.Email),
		Website:  comment.Website,
		Content:  comment.Content,
		CreateAt: comment.CreateAt,
	}
}

func (c *Comment) ToAPIModel() *api.Comment {
	return &api.Comment{
		Nickname:  c.Nickname,
		Avatar:    config.GetGravatarCDN() + c.EmailMD5,
		Website:   c.Website,
		Content:   c.Content,
		CommentAt: "",
	}
}

func (c *Comment) ToAPIModelWithAtUser(at string) *api.Comment {
	return &api.Comment{
		ID:        c.ID,
		Nickname:  c.Nickname,
		Avatar:    config.GetGravatarCDN() + c.EmailMD5,
		Website:   c.Website,
		Content:   c.Content,
		ReplyUser: at,
		CommentAt: "",
	}
}
