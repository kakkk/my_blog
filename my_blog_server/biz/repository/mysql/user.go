package mysql

import (
	"gorm.io/gorm"

	"my_blog/biz/infra/repository/model"
)

func SelectUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	user := &model.User{}
	err := db.Model(&model.User{}).
		Where("username = ?", username).
		First(user).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	return user, nil
}

func SelectUserByID(db *gorm.DB, userID int64) (*model.User, error) {
	user := &model.User{}
	err := db.Model(&model.User{}).
		Where("id = ?", userID).
		First(user).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	return user, nil
}

func MSelectUserByIDs(db *gorm.DB, userIDs []int64) (map[int64]*model.User, error) {
	var users []*model.User
	err := db.Model(&model.User{}).
		Where("id in (?)", userIDs).
		First(&users).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	result := make(map[int64]*model.User, len(userIDs))
	for _, user := range users {
		result[user.ID] = user
	}
	return result, nil
}
