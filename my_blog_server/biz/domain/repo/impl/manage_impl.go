package impl

import (
	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
)

type ManageRepoImpl struct{}

func (m ManageRepoImpl) SelectUserByUsername(db *gorm.DB, username string) (*dto.User, error) {
	user, err := persistence.SelectUserByUsername(db, username)
	if err != nil {
		return nil, err
	}
	return dto.NewUserByModel(user), nil
}

func (m ManageRepoImpl) MSelectUserByIDs(db *gorm.DB, userIDs []int64) (map[int64]*dto.User, error) {
	userMap, err := persistence.MSelectUserByIDs(db, userIDs)
	if err != nil {
		return nil, err
	}
	return dto.NewUserMapByModelMap(userMap), nil
}

func (m ManageRepoImpl) SelectUserByID(db *gorm.DB, userID int64) (*dto.User, error) {
	user, err := persistence.SelectUserByID(db, userID)
	if err != nil {
		return nil, err
	}
	return dto.NewUserByModel(user), nil
}
