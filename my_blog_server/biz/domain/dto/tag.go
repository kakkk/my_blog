package dto

import (
	"github.com/spf13/cast"

	"my_blog/biz/infra/repository/model"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/page"
)

type Tag struct {
	ID      int64
	TagName string
	Count   int64
}

func NewTagByModel(t *model.Tag) *Tag {
	return &Tag{
		ID:      t.ID,
		TagName: t.TagName,
	}
}

func NewTagsByModelMap(tags map[int64]*model.Tag) []*Tag {
	result := make([]*Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, NewTagByModel(tag))
	}
	return result
}

func NewTagsByModel(tags []*model.Tag) []*Tag {
	result := make([]*Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, NewTagByModel(tag))
	}
	return result
}

type Tags []*Tag

func (t Tags) ToStringList() []string {
	result := make([]string, 0, len(t))
	for _, tag := range t {
		result = append(result, tag.TagName)
	}
	return result
}

func (t Tags) ToTagList() []*api.TagListItem {
	tagList := make([]*api.TagListItem, 0, len(t))
	for _, tag := range t {
		tagList = append(tagList, &api.TagListItem{
			ID:    tag.ID,
			Name:  tag.TagName,
			Count: tag.Count,
		})
	}
	return tagList
}

func (t Tags) ToTermListItem() []*page.TermListItem {
	result := make([]*page.TermListItem, 0, len(t))
	for _, item := range t {
		result = append(result, &page.TermListItem{
			Name:  item.TagName,
			Count: cast.ToString(item.Count),
			Slug:  item.TagName,
		})
	}
	return result
}
