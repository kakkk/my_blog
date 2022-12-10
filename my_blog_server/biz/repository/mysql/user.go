package mysql

import (
	"context"
	"fmt"

	"my_blog/biz/common/consts"
	"my_blog/biz/entity"

	"gorm.io/gorm"
)

func SelectUserByUsername(ctx context.Context, db *gorm.DB, username string) (*entity.User, error) {
	user := &entity.User{}
	err := db.Model(&entity.User{}).
		Where("username = ?", username).
		First(user).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("db error: [%v]", err)
	}
	return user, nil
}

func SelectUserByID(ctx context.Context, db *gorm.DB, userID int64) (*entity.User, error) {
	user := &entity.User{}
	err := db.Model(&entity.User{}).
		Where("id = ?", userID).
		First(user).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, consts.ErrRecordNotFound
		}
		return nil, fmt.Errorf("db error: [%v]", err)
	}
	return user, nil
}
