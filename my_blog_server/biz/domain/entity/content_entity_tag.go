package entity

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/idgen"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"
)

type Tag struct {
	ID      int64
	TagName string
}

func (t *Tag) Create(ctx context.Context) error {
	t.ID = idgen.GenID()
	_, err := repo.GetContentRepo().CreateTag(mysql.GetDB(ctx), &dto.Tag{
		ID:      t.ID,
		TagName: t.TagName,
	})
	return err
}

func (t *Tag) Update(ctx context.Context) error {
	return repo.GetContentRepo().UpdateTag(mysql.GetDB(ctx), &dto.Tag{
		ID:      t.ID,
		TagName: t.TagName,
	})
}

func (t *Tag) Delete(ctx context.Context) error {
	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		err := repo.GetContentRepo().DeleteTagByID(tx, t.ID)
		if err != nil {
			return err
		}
		err = repo.GetContentRepo().DeleteArticleTagRelationByTagID(tx, t.ID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func NewTagByDTO(tag *dto.Tag) *Tag {
	return &Tag{
		ID:      tag.ID,
		TagName: tag.TagName,
	}
}

func NewTagByID(id int64) *Tag {
	return &Tag{
		ID: id,
	}
}

type Tags []*Tag

func (t Tags) ToStringList() []string {
	result := make([]string, 0, len(t))
	for _, tag := range t {
		result = append(result, tag.TagName)
	}
	return result
}

func (t Tags) IDs() []int64 {
	result := make([]int64, 0, len(t))
	for _, tag := range t {
		result = append(result, tag.ID)
	}
	return result
}

func NewTagsByDTO(tags []*dto.Tag) Tags {
	result := make([]*Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, &Tag{
			ID:      tag.ID,
			TagName: tag.TagName,
		})
	}
	return result
}

func NewTagsByModelMap(tags map[int64]*model.Tag) []*Tag {
	result := make([]*Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, &Tag{
			ID:      tag.ID,
			TagName: tag.TagName,
		})
	}
	return result
}
