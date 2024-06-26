package model

import (
	"time"

	"my_blog/biz/hertz_gen/blog/common"
)

type ExtraInfo struct {
	InfoType common.ExtraInfo `gorm:"column:id"`
	Value    string           `gorm:"column:value"`
	UpdateAt time.Time        `gorm:"column:update_at"`
}

func (ExtraInfo) TableName() string {
	return "extra_info"
}
