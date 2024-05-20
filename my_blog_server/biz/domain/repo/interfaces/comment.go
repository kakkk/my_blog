package interfaces

import (
	"context"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
)

type CommentRepo interface {
	Create(db *gorm.DB, comment *dto.Comment) error
	GetListByPage(db *gorm.DB, page *int32, size *int32) ([]*dto.Comment, error)
	GetCount(db *gorm.DB) (int64, error)
	GetCommentByID(db *gorm.DB, id int64) (*dto.Comment, error)
	MGetCommentsByID(db *gorm.DB, ids []int64) (map[int64]*dto.Comment, error)
	DeleteByID(db *gorm.DB, id int64) error
	Cache() CommentCache
}

type CommentCache interface {
	Get(ctx context.Context, id int64) (*dto.Comment, error)
	MGet(ctx context.Context, ids []int64) map[int64]*dto.Comment
	GetArticleComments(ctx context.Context, id int64) []*dto.Comment
	RefreshArticleComments(ctx context.Context, id int64)
}
