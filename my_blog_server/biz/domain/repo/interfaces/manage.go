package interfaces

import (
	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
)

type ManageRepo interface {
	SelectUserByUsername(db *gorm.DB, username string) (*dto.User, error)
	MSelectUserByIDs(db *gorm.DB, userIDs []int64) (map[int64]*dto.User, error)
	SelectUserByID(db *gorm.DB, userID int64) (*dto.User, error)
}
