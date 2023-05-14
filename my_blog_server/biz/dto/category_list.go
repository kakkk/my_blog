package dto

import (
	"github.com/spf13/cast"

	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/page"
)

type CategoryListItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Count int64  `json:"count"`
}

type CategoryList []*CategoryListItem

func (c *CategoryList) ToAPICategoryListModel() []*api.CategoryListItem {
	result := make([]*api.CategoryListItem, 0, len(*c))
	for _, item := range *c {
		result = append(result, &api.CategoryListItem{
			ID:    item.ID,
			Name:  item.Name,
			Slug:  item.Slug,
			Count: item.Count,
		})
	}
	return result
}

func (c *CategoryList) ToPageCategoryListModel() []*page.TermListItem {
	result := make([]*page.TermListItem, 0, len(*c))
	for _, item := range *c {
		result = append(result, &page.TermListItem{
			Name:  item.Name,
			Count: cast.ToString(item.Count),
			Slug:  item.Slug,
		})
	}
	return result
}
