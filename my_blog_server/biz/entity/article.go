package entity

import (
	"encoding/json"
	"time"

	"my_blog/biz/model/blog/common"
)

type Article struct {
	ID          int64                `gorm:"column:id"`
	Title       string               `gorm:"column:title"`
	Content     string               `gorm:"column:content"`
	ArticleType common.ArticleType   `gorm:"column:article_type"`
	Status      common.ArticleStatus `gorm:"column:article_status"`
	Extra       string               `gorm:"column:extra"`
	CreateUser  int64                `gorm:"column:create_user"`
	PV          int64                `gorm:"column:pv"`
	CreateAt    time.Time            `gorm:"column:create_at"`
	UpdateAt    time.Time            `gorm:"column:update_at"`
	PublishAt   *time.Time           `gorm:"column:publish_at"`
	DeleteFlag  common.DeleteFlag    `gorm:"column:delete_flag"`
}

func (*Article) TableName() string {
	return "article"
}

func (a *Article) Serialize() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a *Article) Deserialize(str string) (*Article, error) {
	article := &Article{}
	err := json.Unmarshal([]byte(str), article)
	if err != nil {
		return nil, err
	}
	return article, nil
}
