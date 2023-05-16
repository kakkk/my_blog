package entity

import (
	"time"

	"my_blog/biz/model/blog/common"
)

type Category struct {
	ID           int64             `gorm:"column:id"`
	CategoryName string            `gorm:"column:category_name"`
	Slug         string            `gorm:"column:slug"`
	DeleteFlag   common.DeleteFlag `gorm:"column:delete_flag"`
	UpdateAt     time.Time         `gorm:"column:update_at"`
}

func (*Category) TableName() string {
	return "categories"
}

type ArticleCategory struct {
	PostID     int64             `gorm:"column:article_id"`
	CategoryID int64             `gorm:"column:category_id"`
	DeleteFlag common.DeleteFlag `gorm:"column:delete_flag"`
	PublishAt  *time.Time        `gorm:"publish_at"`
}

func (ArticleCategory) TableName() string {
	return "article_category"
}
