package persistence

import (
	"time"

	"gorm.io/gorm"

	"my_blog/biz/consts"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"
)

func CreateComment(db *gorm.DB, comment *model.Comment) (*model.Comment, error) {
	comment.CreateAt = time.Now()
	comment.UpdateAt = time.Now()
	err := db.Model(&model.Comment{}).Create(comment).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return comment, nil
}

func SelectApprovalCommentByID(db *gorm.DB, id int64) (*model.Comment, error) {
	comment := &model.Comment{}
	err := db.Model(&model.Comment{}).
		Where("id = ?", id).
		Where("status = ?", common.CommentStatus_Approved).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(comment).
		Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return comment, nil
}

func SelectCommentByID(db *gorm.DB, id int64) (*model.Comment, error) {
	comment := &model.Comment{}
	err := db.Model(&model.Comment{}).
		Where("id = ?", id).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(comment).
		Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return comment, nil
}

func SelectApprovalCommentsByIDs(db *gorm.DB, ids []int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := db.Model(&model.Comment{}).
		Where("id in ?", ids).
		Where("status = ?", common.CommentStatus_Approved).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&comments).
		Error
	if err != nil {
		return nil, mysql.ParseError(err)
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
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&ids).
		Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	if len(ids) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return ids, nil
}

func SelectCommentByPage(db *gorm.DB, page *int32, size *int32) ([]*model.Comment, error) {
	var result []*model.Comment
	offset, limit := 0, mysql.DefaultPageLimit
	if size != nil {
		limit = int(*size)
	}
	if page != nil {
		if *page != 0 {
			offset = (int(*page) - 1) * limit
		}
	}
	err := db.Model(&model.Comment{}).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Limit(limit).
		Offset(offset).
		Order("id desc").
		Find(&result).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return result, nil
}

func SelectCommentsByIDs(db *gorm.DB, ids []int64) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := db.Model(&model.Comment{}).
		Where("id in ?", ids).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&comments).
		Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	if len(comments) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return comments, nil
}

func DeleteCommentByID(db *gorm.DB, id int64) error {
	err := db.Model(&model.Comment{}).
		Where("id = ?", id).
		Update("delete_flag", common.DeleteFlag_Delete).
		Error
	if err != nil {
		return mysql.ParseError(err)
	}
	return nil
}

func SelectCommentCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&model.Comment{}).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Count(&count).Error
	if err != nil {
		return 0, mysql.ParseError(err)
	}
	return count, nil
}
