package dto

import (
	"github.com/spf13/cast"

	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/hertz_gen/blog/page"
	"my_blog/biz/infra/repository/model"
)

type Category struct {
	ID           int64  `json:"id"`
	CategoryName string `json:"category_name"`
	Slug         string `json:"slug"`
	Count        int64  `json:"count"`
}

func NewCategoryByModel(c *model.Category) *Category {
	return &Category{
		ID:           c.ID,
		CategoryName: c.CategoryName,
		Slug:         c.Slug,
		Count:        0,
	}
}

type Categories []*Category

func NewCategoriesByModelMap(categories map[int64]*model.Category) []*Category {
	result := make([]*Category, 0, len(categories))
	for _, item := range categories {
		if item != nil {
			result = append(result, NewCategoryByModel(item))
		}
	}
	return result
}

func NewCategoriesByModelList(categories []*model.Category) []*Category {
	result := make([]*Category, 0, len(categories))
	for _, item := range categories {
		if item != nil {
			result = append(result, NewCategoryByModel(item))
		}
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

func (c Categories) ToTermListItem() []*page.TermListItem {
	result := make([]*page.TermListItem, 0, len(c))
	for _, item := range c {
		result = append(result, &page.TermListItem{
			Name:  item.CategoryName,
			Count: cast.ToString(item.Count),
			Slug:  item.Slug,
		})
	}
	return result
}

func (c Categories) ToStringList() []string {
	result := make([]string, 0, len(c))
	for _, category := range c {
		result = append(result, category.CategoryName)
	}
	return result
}
