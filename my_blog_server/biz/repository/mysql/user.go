package mysql

import (
	"my_blog/biz/entity"

	"gorm.io/gorm"
)

func SelectUserByUsername(db *gorm.DB, username string) (*entity.User, error) {
	user := &entity.User{}
	err := db.Model(&entity.User{}).
		Where("username = ?", username).
		First(user).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	return user, nil
}

func SelectUserByID(db *gorm.DB, userID int64) (*entity.User, error) {
	user := &entity.User{}
	err := db.Model(&entity.User{}).
		Where("id = ?", userID).
		First(user).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	return user, nil
}
