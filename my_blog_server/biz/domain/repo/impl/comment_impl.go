package impl

import (
	"context"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/interfaces"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/hertz_gen/blog/common"
)

type CommentRepoImpl struct {
	cache CommentCacheImpl
}

func (c CommentRepoImpl) Create(db *gorm.DB, comment *dto.Comment) error {
	_, err := persistence.CreateComment(db, comment.ToModel())
	return err
}

func (c CommentRepoImpl) GetListByPage(db *gorm.DB, page *int32, size *int32) ([]*dto.Comment, error) {
	comments, err := persistence.SelectCommentByPage(db, page, size)
	if err != nil {
		return nil, err
	}
	result := dto.NewCommentListByModel(comments)
	return result, nil
}

func (c CommentRepoImpl) GetCount(db *gorm.DB) (int64, error) {
	return persistence.SelectCommentCount(db)
}

func (c CommentRepoImpl) GetCommentByID(db *gorm.DB, id int64) (*dto.Comment, error) {
	comment, err := persistence.SelectCommentByID(db, id)
	if err != nil {
		return nil, err
	}
	return dto.NewCommentByModel(comment), nil
}

func (c CommentRepoImpl) MGetCommentsByID(db *gorm.DB, ids []int64) (map[int64]*dto.Comment, error) {
	comments, err := persistence.SelectCommentsByIDs(db, ids)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*dto.Comment)
	for _, comment := range comments {
		result[comment.ID] = dto.NewCommentByModel(comment)
	}
	return result, nil
}

func (c CommentRepoImpl) DeleteByID(db *gorm.DB, id int64) error {
	return persistence.DeleteCommentByID(db, id)
}

func (c CommentRepoImpl) UpdateCommentStatusByID(db *gorm.DB, id int64, status common.CommentStatus) error {
	return persistence.UpdateCommentStatusByID(db, id, status)
}

func (c CommentRepoImpl) Cache() interfaces.CommentCache {
	return c.cache
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

func (c CommentCacheImpl) RefreshArticleComments(ctx context.Context, id int64) {
	articleCommentsCachex.Refresh(ctx, id)
}
