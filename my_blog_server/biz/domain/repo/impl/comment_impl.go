package impl

import (
	"context"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/interfaces"
	"my_blog/biz/domain/repo/persistence"
)

type CommentRepoImpl struct{}

func (c CommentRepoImpl) Create(db *gorm.DB, comment *dto.Comment) error {
	_, err := persistence.CreateComment(db, comment.ToModel())
	return err
}

func (c CommentRepoImpl) Cache() interfaces.CommentCache {
	//TODO implement me
	panic("implement me")
}

type CommentCacheImpl struct{}

func (c CommentCacheImpl) Get(ctx context.Context, id int64) (*dto.Comment, error) {
	return commentCachex.Get(ctx, id)
}

func (c CommentCacheImpl) MGet(ctx context.Context, ids []int64) map[int64]*dto.Comment {
	return commentCachex.MGet(ctx, ids)
}

func (c CommentCacheImpl) GetArticleComments(ctx context.Context, id int64) []*dto.Comment {
	return articleCommentsCachex.Get(ctx, id)
}
