package mysql

import (
	"time"

	"gorm.io/gorm"

	"my_blog/biz/common/consts"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/common"
)

func CreateComment(db *gorm.DB, comment *entity.Comment) (*entity.Comment, error) {
	comment.CreateAt = time.Now()
	comment.UpdateAt = time.Now()
	err := db.Model(&entity.Comment{}).Create(comment).Error
	if err != nil {
		return nil, parseError(err)
	}
	return comment, nil
}

func SelectApprovalCommentByID(db *gorm.DB, id int64) (*entity.Comment, error) {
	comment := &entity.Comment{}
	err := db.Model(&entity.Comment{}).
		Where("id = ?", id).
		Where("status = ?", common.CommentStatus_Approved).
		First(comment).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	return comment, nil
}

func SelectApprovalCommentsByIDs(db *gorm.DB, ids []int64) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	err := db.Model(&entity.Comment{}).
		Where("id in ?", ids).
		Where("status = ?", common.CommentStatus_Approved).
		Find(&comments).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	if len(comments) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return comments, nil
}

func SelectCommentIDsByArticleID(db *gorm.DB, id int64) ([]int64, error) {
	var ids []int64
	err := db.Model(&entity.Comment{}).
		Select("id").
		Where("article_id = ?", id).
		Where("status = ?", common.CommentStatus_Approved).
		Find(&ids).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	if len(ids) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return ids, nil
}
