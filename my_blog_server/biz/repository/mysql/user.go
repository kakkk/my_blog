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

func MSelectUserByIDs(db *gorm.DB, userIDs []int64) (map[int64]*entity.User, error) {
	var users []*entity.User
	err := db.Model(&entity.User{}).
		Where("id in (?)", userIDs).
		First(&users).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	result := make(map[int64]*entity.User, len(userIDs))
	for _, user := range users {
		result[user.ID] = user
	}
	return result, nil
}
