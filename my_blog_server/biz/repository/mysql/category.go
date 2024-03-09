package mysql

import (
	"encoding/json"
	"fmt"
	"time"

	"my_blog/biz/infra/misc"
	model2 "my_blog/biz/infra/repository/model"
	"my_blog/biz/model/blog/common"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateCategory(db *gorm.DB, category *model2.Category) (*model2.Category, error) {
	category.UpdateAt = time.Now()
	err := db.Model(&model2.Category{}).Create(category).Error
	if err != nil {
		return nil, parseError(err)
	}
	return category, nil
}

func SelectCategoryByID(db *gorm.DB, id int64) (*model2.Category, error) {
	got := &model2.Category{}
	err := db.Model(&model2.Category{}).
		Where("id = ?", id).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&got).Error
	if err != nil {
		return nil, parseError(err)
	}
	return got, nil
}

func MSelectCategoryByIDs(db *gorm.DB, ids []int64) (map[int64]*model2.Category, error) {
	got := make([]*model2.Category, 0, len(ids))
	err := db.Model(&model2.Category{}).
		Where("id in (?)", ids).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&got).Error
	if err != nil {
		return nil, parseError(err)
	}
	result := make(map[int64]*model2.Category, len(ids))
	for _, category := range got {
		result[category.ID] = category
	}
	return result, nil
}

func UpdateCategoryByID(db *gorm.DB, categoryID int64, category *model2.Category) error {
	err := db.Model(&model2.Category{}).
		Where("id = ?", categoryID).
		Updates(map[string]any{
			"category_name": category.CategoryName,
			"slug":          category.Slug,
			"update_at":     time.Now(),
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func DeleteCategoryByID(db *gorm.DB, categoryID int64) error {
	err := db.Model(&model2.Category{}).
		Where("id = ?", categoryID).
		Updates(map[string]any{
			"delete_flag": common.DeleteFlag_Delete,
			"update_at":   time.Now(),
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func SelectCategoryOrder(db *gorm.DB) ([]int64, error) {
	value := ""
	err := db.Model(&model2.ExtraInfo{}).
		Select("value").
		Where("id = ?", common.ExtraInfo_CategoryOrder).
		First(&value).Error
	if err != nil {
		return nil, parseError(err)
	}
	var result []int64
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, fmt.Errorf("json error:[%v]", err)
	}
	return result, nil
}

func UpdateCategoryOrder(db *gorm.DB, order []int64) error {
	value, err := json.Marshal(&order)
	if err != nil {
		return fmt.Errorf("json error:[%v]", err)
	}
	err = db.Model(&model2.ExtraInfo{}).
		Where("id = ?", common.ExtraInfo_CategoryOrder).
		Updates(map[string]any{
			"value":     string(value),
			"update_at": time.Now(),
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func MSelectCategoryArticleCountByCategoryIDs(db *gorm.DB, categoryIDs []int64, withPublish bool) (map[int64]int64, error) {
	type result struct {
		CategoryID int64 `gorm:"column:category_id"`
		Count      int64 `gorm:"column:count"`
	}
	res := make(map[int64]int64, len(categoryIDs))
	for _, id := range categoryIDs {
		res[id] = 0
	}
	var resultFromDB []result
	query := db.Model(&model2.ArticleCategory{}).
		Select("category_id, count(1) as count").
		Where("category_id in (?)", categoryIDs).
		Where("delete_flag = ?", common.DeleteFlag_Exist)
	if withPublish {
		query.Where("publish_at is not null")
	}
	err := query.Group("category_id").Find(&resultFromDB).Error
	if err != nil {
		return res, parseError(err)
	}

	for _, item := range resultFromDB {
		res[item.CategoryID] = item.Count
	}

	return res, nil
}

func UpsertArticleCategoryRelation(db *gorm.DB, articleCategories []*model2.ArticleCategory) error {
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "article_id"}, {Name: "category_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"delete_flag"}),
	}).Create(&articleCategories).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func SelectDefaultCategoryID(db *gorm.DB) (int64, error) {
	value := int64(0)
	err := db.Model(&model2.ExtraInfo{}).
		Select("value").
		Where("id = ?", common.ExtraInfo_DefaultCategory).
		First(&value).Error
	if err != nil {
		return 0, parseError(err)
	}

	return value, nil
}

func SelectCategoryIDsByArticleID(db *gorm.DB, articleID int64) ([]int64, error) {
	var result []int64
	err := db.Model(&model2.ArticleCategory{}).
		Select("category_id").
		Where("article_id = ?", articleID).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&result).Error
	if err != nil {
		return nil, parseError(err)
	}
	return result, nil
}

func DeleteArticleCategoryRelationByArticleID(db *gorm.DB, articleID int64) error {
	err := db.Model(&model2.ArticleCategory{}).
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

func SelectArticleIDsByCategoryIDs(db *gorm.DB, categoryIDs []int64) ([]int64, error) {
	if len(categoryIDs) == 0 {
		return []int64{}, nil
	}
	var results []int64
	err := db.Model(&model2.ArticleCategory{}).
		Select("article_id").
		Where("category_id in (?)", categoryIDs).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&results).Error
	if err != nil {
		return nil, parseError(err)
	}
	return misc.SliceDeduplicate[int64](results), nil
}

func MSelectCategoryIDsByArticleIDs(db *gorm.DB, articleIDs []int64) (map[int64][]int64, error) {
	result := map[int64][]int64{}
	for _, id := range articleIDs {
		categoryIDs, err := SelectCategoryIDsByArticleID(db, id)
		if err != nil {
			return nil, parseError(err)
		}
		result[id] = categoryIDs
	}
	return result, nil
}

func MSelectCategoryByNames(db *gorm.DB, names []string) (map[string]*model2.Category, error) {
	got := make([]*model2.Category, 0, len(names))
	err := db.Model(&model2.Category{}).
		Where("category_name in (?)", names).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&got).Error
	if err != nil {
		return nil, parseError(err)
	}
	result := make(map[string]*model2.Category, len(names))
	for _, category := range got {
		result[category.CategoryName] = category
	}
	return result, nil
}

func SelectCategoryIDBySlug(db *gorm.DB, slug string) (int64, error) {
	var id int64
	err := db.Model(&model2.Category{}).
		Select("id").
		Where("slug = ?", slug).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(&id).
		Error
	if err != nil {
		return 0, parseError(err)
	}
	return id, nil
}

func UpdateArticleCategoryUpdateAtByArticleID(db *gorm.DB, id int64, publishAt *time.Time) error {
	err := db.Model(&model2.ArticleCategory{}).
		Where("article_id = ?", id).
		Update("publish_at", publishAt).Error
	return parseError(err)
}
