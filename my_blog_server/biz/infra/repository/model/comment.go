package model

import (
	"time"

	"my_blog/biz/hertz_gen/blog/common"
)

type Comment struct {
	ID         int64                `gorm:"column:id"`
	ArticleID  int64                `gorm:"column:article_id"`
	ReplyID    int64                `gorm:"column:reply_id"`
	ParentID   int64                `gorm:"column:parent_id"`
	Nickname   string               `gorm:"column:nickname"`
	Email      string               `gorm:"column:email"`
	Website    string               `gorm:"column:website"`
	Content    string               `gorm:"column:content"`
	Status     common.CommentStatus `gorm:"column:status"`
	CreateAt   time.Time            `gorm:"column:create_at"`
	UpdateAt   time.Time            `gorm:"column:update_at"`
	DeleteFlag common.DeleteFlag    `gorm:"column:delete_flag"`
}

func (Comment) TableName() string {
	return "comment"
}
