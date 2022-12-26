package entity

import (
	"time"

	"my_blog/biz/model/blog/common"
)

type Tag struct {
	ID         int64             `gorm:"column:id"`
	TagName    string            `gorm:"column:tag_name"`
	UpdateAt   time.Time         `gorm:"column:update_at"`
	DeleteFlag common.DeleteFlag `gorm:"column:delete_flag"`
}

func (Tag) TableName() string {
	return "tags"
}

type ArticleTag struct {
	PostID     int64             `gorm:"column:article_id"`
	TagID      int64             `gorm:"column:tag_id"`
	DeleteFlag common.DeleteFlag `gorm:"column:delete_flag"`
}

func (ArticleTag) TableName() string {
	return "article_tag"
}
