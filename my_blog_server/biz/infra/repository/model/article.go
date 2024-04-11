package model

import (
	"time"

	"my_blog/biz/hertz_gen/blog/common"
)

type Article struct {
	ID          int64                `gorm:"column:id"`
	Title       string               `gorm:"column:title"`
	Content     string               `gorm:"column:content"`
	ArticleType common.ArticleType   `gorm:"column:article_type"`
	Status      common.ArticleStatus `gorm:"column:article_status"`
	Extra       string               `gorm:"column:extra"`
	CreateUser  int64                `gorm:"column:create_user"`
	UV          int64                `gorm:"column:uv"`
	CreateAt    time.Time            `gorm:"column:create_at"`
	UpdateAt    time.Time            `gorm:"column:update_at"`
	PublishAt   *time.Time           `gorm:"column:publish_at"`
	DeleteFlag  common.DeleteFlag    `gorm:"column:delete_flag"`
}

func (*Article) TableName() string {
	return "article"
}
