package mysql

import (
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/common"

	"gorm.io/gorm"
)

const (
	defaultPageLimit = 10
)

func CreateTag(db *gorm.DB, tag *entity.Tag) (*entity.Tag, error) {
	err := db.Model(&entity.Tag{}).Create(tag).Error
	if err != nil {
		return nil, parseError(err)
	}
	return tag, nil
}

func UpdateTagByID(db *gorm.DB, id int64, tag *entity.Tag) error {
	err := db.Model(&entity.Tag{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"name": tag.TagName,
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func DeleteTagByID(db *gorm.DB, id int64) error {
	err := db.Model(&entity.Tag{}).
		Where("id = ?", id).
		Updates(
			map[string]any{
				"delete_flag": common.DeleteFlag_Delete,
			},
		).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func GetTagListByPage(db *gorm.DB, keyword *string, page *int32, size *int32) ([]*entity.Tag, error) {
	var result []*entity.Tag
	tx := db.Model(&entity.Tag{})
	if keyword != nil {
		tx.Where("tag_name like ?", *keyword+"%")
	}
	tx.Where("delete_flag = ?", common.DeleteFlag_Exist)

	offset, limit := 0, defaultPageLimit
	if size != nil {
		limit = int(*size)
	}
	if page != nil {
		if *page != 0 {
			offset = (int(*page) - 1) * limit
		}
	}
	err := tx.Limit(limit).
		Offset(offset).
		Order("update_at desc").
		Find(&result).Error
	if err != nil {
		return nil, parseError(err)
	}
	return result, nil
}

func MGetTagArticleCountByTagIDs(db *gorm.DB, tagIDs []int64) (map[int64]int64, error) {
	type result struct {
		TagID int64 `gorm:"column:tag_id"`
		Count int64 `gorm:"column:count"`
	}
	var resultFromDB []result
	err := db.Model(&entity.ArticleTag{}).
		Select("tag_id, count(1) as count").
		Where("tag_id in (?)", tagIDs).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Group("tag_id").
		Find(&resultFromDB).Error
	if err != nil {
		return nil, parseError(err)
	}

	res := make(map[int64]int64, len(resultFromDB))
	for _, id := range tagIDs {
		res[id] = 0
	}
	for _, item := range resultFromDB {
		res[item.TagID] = item.Count
	}

	return res, nil
}

func GetAllTagCount(db *gorm.DB) (int64, error) {
	count := int64(0)
	err := db.Model(&entity.Tag{}).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Count(&count).Error
	if err != nil {
		return 0, parseError(err)
	}
	return count, nil
}

func DeleteTagArticleByTagID(db *gorm.DB, tagID int64) error {
	err := db.Model(&entity.ArticleTag{}).
		Where("tag_id = ?", tagID).
		Updates(
			map[string]any{
				"delete_flag": common.DeleteFlag_Delete,
			},
		).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}
