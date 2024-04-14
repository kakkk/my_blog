package interfaces

import (
	"context"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
)

type CommentRepo interface {
	Create(db *gorm.DB, comment *dto.Comment) error
	Cache() CommentCache
}

type CommentCache interface {
	Get(ctx context.Context, id int64) (*dto.Comment, error)
	MGet(ctx context.Context, ids []int64) map[int64]*dto.Comment
	GetArticleComments(ctx context.Context, id int64) []*dto.Comment
	RefreshArticleComments(ctx context.Context, id int64)
}
