package mysql

import (
	"fmt"
	"time"

	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateTag(db *gorm.DB, tag *entity.Tag) (*entity.Tag, error) {
	tag.UpdateAt = time.Now()
	err := db.Model(&entity.Tag{}).Create(tag).Error
	if err != nil {
		return nil, parseError(err)
	}
	return tag, nil
}

func BatchCreateTags(db *gorm.DB, tags []*entity.Tag) ([]*entity.Tag, error) {
	if len(tags) == 0 {
		return nil, nil
	}
	now := time.Now()
	for i := range tags {
		tags[i].UpdateAt = now
	}
	err := db.Model(&entity.Tag{}).Create(&tags).Error
	if err != nil {
		return nil, parseError(err)
	}
	return tags, nil
}

func UpdateTagByID(db *gorm.DB, id int64, tag *entity.Tag) error {
	err := db.Model(&entity.Tag{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"name":      tag.TagName,
			"update_at": time.Now(),
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
				"update_at":   time.Now(),
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

func DeleteArticleTagRelationByTagID(db *gorm.DB, tagID int64) error {
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

func MSelectTagByName(db *gorm.DB, names []string) (map[string]*entity.Tag, error) {
	tags := make([]*entity.Tag, 0, len(names))
	err := db.Model(&entity.Tag{}).
		Where("tag_name in (?)", names).
		Find(&tags).Error
	if err != nil {
		return nil, parseError(err)
	}
	result := make(map[string]*entity.Tag, len(tags))
	for _, tag := range tags {
		result[tag.TagName] = tag
	}
	return result, nil
}

func RestoreTagByIDs(db *gorm.DB, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	err := db.Model(&entity.Tag{}).
		Where("id in (?)", ids).
		Updates(
			map[string]any{
				"delete_flag": common.DeleteFlag_Exist,
				"update_at":   time.Now(),
			},
		).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func UpsertArticleTagRelation(db *gorm.DB, articleTags []*entity.ArticleTag) error {
	if len(articleTags) == 0 {
		return nil
	}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "article_id"}, {Name: "tag_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"delete_flag"}),
	}).Create(&articleTags).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func DeleteArticleTagRelationByArticleID(db *gorm.DB, articleID int64) error {
	err := db.Model(&entity.ArticleTag{}).
		Where("article_id = ?", articleID).
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

func RefreshTagUpdateTimeByIDs(db *gorm.DB, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	err := db.Model(&entity.Tag{}).
		Where("id in (?)", ids).
		Updates(
			map[string]any{
				"update_at": time.Now(),
			},
		).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func SelectTagIDsByArticleID(db *gorm.DB, articleID int64) ([]int64, error) {
	var result []int64
	err := db.Model(&entity.ArticleTag{}).
		Select("tag_id").
		Where("article_id = ?", articleID).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&result).Error
	if err != nil {
		return nil, parseError(err)
	}
	return result, nil
}

func MSelectTagByID(db *gorm.DB, ids []int64) (map[int64]*entity.Tag, error) {
	tags := make([]*entity.Tag, 0, len(ids))
	err := db.Model(&entity.Tag{}).
		Where("id in (?)", ids).
		Find(&tags).Error
	if err != nil {
		return nil, parseError(err)
	}
	result := make(map[int64]*entity.Tag, len(tags))
	for _, tag := range tags {
		result[tag.ID] = tag
	}
	return result, nil
}

func SelectArticleIDsByTagIDs(db *gorm.DB, ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return []int64{}, nil
	}
	var results []int64
	err := db.Model(&entity.ArticleTag{}).
		Select("article_id").
		Where("tag_id in (?)", ids).
		Find(&results).Error
	if err != nil {
		return nil, parseError(err)
	}
	return utils.SliceDeduplicate[int64](results), nil
}

func SelectTagListByArticleID(db *gorm.DB, articleID int64) ([]string, error) {
	// 获取标签
	tagIDs, err := SelectTagIDsByArticleID(db, articleID)
	if err != nil {
		return nil, fmt.Errorf("select tag_id error:[%v]", err)
	}
	tagMap, err := MSelectTagByID(db, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("select tag error:[%v]", err)
	}
	var tagList []string
	for _, tag := range tagMap {
		tagList = append(tagList, tag.TagName)
	}
	return tagList, nil
}

func UpdateArticleTagUpdateAtByArticleID(db *gorm.DB, id int64, publishAt *time.Time) error {
	err := db.Model(&entity.ArticleTag{}).
		Where("article_id = ?", id).
		Update("publish_at", publishAt).Error
	return parseError(err)
}
