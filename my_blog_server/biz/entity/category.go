package entity

import (
	"encoding/json"
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

func (Category) TableName() string {
	return "categories"
}

func (a *Category) Serialize() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a *Category) Deserialize(str string) (*Category, error) {
	category := &Category{}
	err := json.Unmarshal([]byte(str), category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

type ArticleCategory struct {
	PostID     int64             `gorm:"column:article_id"`
	CategoryID int64             `gorm:"column:category_id"`
	DeleteFlag common.DeleteFlag `gorm:"column:delete_flag"`
}

func (ArticleCategory) TableName() string {
	return "article_category"
}
