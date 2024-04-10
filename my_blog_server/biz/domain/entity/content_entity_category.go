package entity

import (
	"context"
	"fmt"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/api"
)

type Category struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"category_name"`
	Slug         string `json:"slug"`
	Count        int64  `json:"count"`
}

func (c *Category) Create(ctx context.Context) error {
	tx := mysql.GetDB(ctx).Begin()
	// 创建分类
	category, err := repo.GetContentRepo().CreateCategory(tx, c.ToCategoryDTO())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("create category error:[%w]", err)
	}

	// 更新排序
	order, err := repo.GetContentRepo().SelectCategoryOrder(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("select category order error:[%w]", err)
	}
	order = append(order, category.ID)
	err = repo.GetContentRepo().UpdateCategoryOrder(tx, order)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("update category order error:[%w]", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("commit transaction error:[%v]", err)
	}
	return nil
}

func (c *Category) Update(ctx context.Context) error {
	err := persistence.UpdateCategoryByID(mysql.GetDB(ctx), c.ID, &model.Category{
		CategoryName: c.CategoryName,
		Slug:         c.Slug,
	})
	if err != nil {
		return fmt.Errorf("update category error:[%w]", err)
	}
	return nil
}

func (c *Category) Delete(ctx context.Context) error {
	tx := mysql.GetDB(ctx).Begin()

	// 删除文章<->分类关联关系
	err := persistence.DeleteArticleCategoryRelationByCategoryID(tx, c.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete category article relation error:[%w]", err)
	}

	// 删除分类
	err = persistence.DeleteCategoryByID(tx, c.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete category error:[%w]", err)
	}

	// 更新排序
	beforeOrder, err := persistence.SelectCategoryOrder(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("select category order error:[%w]", err)
	}
	afterOrder := make([]int64, 0, len(beforeOrder)-1)
	for _, id := range beforeOrder {
		if id == c.ID {
			continue
		}
		afterOrder = append(afterOrder, id)
	}
	err = persistence.UpdateCategoryOrder(tx, afterOrder)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("update category order error:[%w]", err)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("commit transaction error:[%w]", err)
	}
	return nil
}

func (c *Category) ToCategoryDTO() *dto.Category {
	return &dto.Category{
		ID:           c.ID,
		CategoryName: c.CategoryName,
		Slug:         c.Slug,
		Count:        c.Count,
	}
}

type Categories []*Category

func (c Categories) ToCategoriesItems() []*api.CategoriesItem {
	result := make([]*api.CategoriesItem, 0, len(c))
	for _, category := range c {
		result = append(result, &api.CategoriesItem{
			ID:   category.ID,
			Name: category.CategoryName,
		})
	}
	return result
}

func (c Categories) ToCategoryListItems() []*api.CategoryListItem {
	result := make([]*api.CategoryListItem, 0, len(c))
	for _, item := range c {
		result = append(result, &api.CategoryListItem{
			ID:    item.ID,
			Name:  item.CategoryName,
			Slug:  item.Slug,
			Count: item.Count,
		})
	}
	return result
}

func NewCategoriesByModel(categories map[int64]*model.Category) []*Category {
	result := make([]*Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, &Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			Slug:         category.Slug,
		})
	}
	return result
}

func NewCategoriesByDTO(categories []*dto.Category) Categories {
	result := make([]*Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, &Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			Slug:         category.Slug,
		})
	}
	return result
}

func NewCategoriesByMap(categories map[int64]*dto.Category) Categories {
	result := make([]*Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, &Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			Slug:         category.Slug,
		})
	}
	return result
}
