package impl

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"my_blog/biz/consts"
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/repository/model"
)

func (c ContentRepoImpl) CreateCategory(db *gorm.DB, category *dto.Category) (*dto.Category, error) {
	categoryModel := &model.Category{
		CategoryName: category.CategoryName,
		Slug:         category.Slug,
	}
	categoryModel, err := persistence.CreateCategory(db, categoryModel)
	if err != nil {
		return nil, err
	}
	category.ID = categoryModel.ID
	return category, nil
}

func (c ContentRepoImpl) SelectCategoryOrder(db *gorm.DB) ([]int64, error) {
	return persistence.SelectCategoryOrder(db)
}

func (c ContentRepoImpl) UpdateCategoryOrder(db *gorm.DB, order []int64) error {
	return persistence.UpdateCategoryOrder(db, order)
}

func (c ContentRepoImpl) GetCategoryList(db *gorm.DB, publish bool) ([]*dto.Category, error) {
	order, err := persistence.SelectCategoryOrder(db)
	if err != nil {
		return nil, fmt.Errorf("select category order error:[%w]", err)
	}
	categories, err := persistence.MSelectCategoryByIDs(db, order)
	if err != nil {
		return nil, fmt.Errorf("select category error:[%w]", err)
	}

	counts, err := persistence.MSelectCategoryArticleCountByCategoryIDs(db, order, publish)
	if err != nil {
		return nil, fmt.Errorf("select article count error:[%w]", err)
	}

	var list []*dto.Category
	for _, id := range order {
		category, ok := categories[id]
		if !ok {
			continue
		}
		count := counts[id]
		list = append(list, &dto.Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			Slug:         category.Slug,
			Count:        count,
		})
	}
	return list, nil
}

func (c ContentRepoImpl) SelectCategoryByIDs(db *gorm.DB, ids []int64) ([]*dto.Category, error) {
	got, err := persistence.MSelectCategoryByIDs(db, ids)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.Category, 0, len(ids))
	for _, id := range ids {
		res, ok := got[id]
		if !ok {
			return nil, fmt.Errorf("%w, id:%v", consts.ErrRecordNotFound, id)
		}
		result = append(result, dto.NewCategoryByModel(res))
	}
	return result, nil
}

func (c ContentRepoImpl) DeleteArticleCategoryRelationByArticleID(db *gorm.DB, articleID int64) error {
	return persistence.DeleteArticleCategoryRelationByArticleID(db, articleID)
}

func (c ContentRepoImpl) UpsertArticleCategoryRelation(db *gorm.DB, articleCategories []*model.ArticleCategory) error {
	return persistence.UpsertArticleCategoryRelation(db, articleCategories)
}

func (c ContentRepoImpl) GetCategoriesByArticleID(db *gorm.DB, articleID int64) ([]*dto.Category, error) {
	// 获取分类
	categoryIDs, err := persistence.SelectCategoryIDsByArticleID(db, articleID)
	if err != nil {
		return nil, fmt.Errorf("select category_id error:[%v]", err)
	}
	if len(categoryIDs) == 0 {
		return make([]*dto.Category, 0), nil
	}
	categoryMap, err := persistence.MSelectCategoryByIDs(db, categoryIDs)
	if err != nil {
		return nil, fmt.Errorf("select category_id error:[%v]", err)
	}
	return dto.NewCategoriesByModelMap(categoryMap), nil
}

func (c ContentRepoImpl) UpdateArticleCategoryUpdateAtByArticleID(db *gorm.DB, id int64, publishAt *time.Time) error {
	return persistence.UpdateArticleCategoryUpdateAtByArticleID(db, id, publishAt)
}

func (c ContentRepoImpl) SelectArticleIDsByCategoryNames(db *gorm.DB, categoryNames []string) ([]int64, error) {
	categoryMap, err := persistence.MSelectCategoryByNames(db, categoryNames)
	if err != nil {
		return nil, err
	}
	var categoryIDs []int64
	for _, category := range categoryMap {
		categoryIDs = append(categoryIDs, category.ID)
	}
	if len(categoryIDs) == 0 {
		return make([]int64, 0), nil
	}
	postIDs, err := persistence.SelectArticleIDsByCategoryIDs(db, categoryIDs)
	if err != nil {
		return nil, err
	}
	return postIDs, nil
}
