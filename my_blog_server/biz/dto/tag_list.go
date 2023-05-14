package dto

import (
	"sort"

	"github.com/spf13/cast"

	"my_blog/biz/entity"
	"my_blog/biz/model/blog/page"
)

type TagListItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TagList []*TagListItem

func (c *TagList) ToPageCategoryListModel() []*page.TermListItem {
	result := make([]*page.TermListItem, 0, len(*c))
	for _, item := range *c {
		result = append(result, &page.TermListItem{
			Name:  item.Name,
			Count: cast.ToString(item.Count),
			Slug:  item.Name,
		})
	}
	return result
}

func NewTagList(tags []*entity.Tag, countMap map[int64]int64) *TagList {
	result := make([]*TagListItem, 0, len(tags))
	for _, tag := range tags {
		result = append(result, &TagListItem{
			ID:    tag.ID,
			Name:  tag.TagName,
			Count: countMap[tag.ID],
		})
	}
	// 排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count >= result[j].Count
	})
	return (*TagList)(&result)
}
