package mysql

import (
	"time"

	"gorm.io/gorm"

	"my_blog/biz/consts"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/model/blog/common"
)

func CreateComment(db *gorm.DB, comment *model.Comment) (*model.Comment, error) {
	comment.CreateAt = time.Now()
	comment.UpdateAt = time.Now()
	err := db.Model(&model.Comment{}).Create(comment).Error
	if err != nil {
		return nil, parseError(err)
	}
	return comment, nil
}

func SelectApprovalCommentByID(db *gorm.DB, id int64) (*model.Comment, error) {
	comment := &model.Comment{}
	err := db.Model(&model.Comment{}).
		Where("id = ?", id).
		Where("status = ?", common.CommentStatus_Approved).
		First(comment).
		Error
	if err != nil {
		return nil, parseError(err)
	}
	return comment, nil
}

func SelectApprovalCommentsByIDs(db *gorm.DB, ids []int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := db.Model(&model.Comment{}).
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
	err := db.Model(&model.Comment{}).
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
