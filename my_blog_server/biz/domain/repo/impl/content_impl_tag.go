package impl

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/repository/model"
)

func (c ContentRepoImpl) GetOrCreateTagsByNames(db *gorm.DB, names []string) ([]*dto.Tag, error) {
	// 获取存在的标签
	tagMap, err := persistence.MSelectTagByName(db, names)
	if err != nil {
		return nil, fmt.Errorf("select tags error:[%v]", err)
	}
	notExistTag := make([]*model.Tag, 0)  // 不存在
	needRestoreTagIDs := make([]int64, 0) // 需要恢复
	tags := make([]*dto.Tag, 0, len(names))
	for _, name := range names {
		tag, ok := tagMap[name]
		if !ok {
			notExistTag = append(notExistTag, &model.Tag{TagName: name})
			continue
		}
		if tag.DeleteFlag == common.DeleteFlag_Delete {
			needRestoreTagIDs = append(needRestoreTagIDs, tag.ID)
		}
		tags = append(tags, dto.NewTagByModel(tag))
	}

	// 创建不存在的标签
	createdTags, err := persistence.BatchCreateTags(db, notExistTag)
	if err != nil {
		return nil, fmt.Errorf("create tags error:[%v]", err)
	}
	for _, tag := range createdTags {
		tags = append(tags, dto.NewTagByModel(tag))
	}

	// 恢复被删除的标签
	err = persistence.RestoreTagByIDs(db, needRestoreTagIDs)
	if err != nil {
		return nil, fmt.Errorf("restore tags error:[%v]", err)
	}

	return tags, nil
}

func (c ContentRepoImpl) DeleteArticleTagRelationByArticleID(db *gorm.DB, articleID int64) error {
	return persistence.DeleteArticleTagRelationByArticleID(db, articleID)
}

func (c ContentRepoImpl) UpsertArticleTagRelation(db *gorm.DB, articleTags []*model.ArticleTag) error {
	return persistence.UpsertArticleTagRelation(db, articleTags)
}

func (c ContentRepoImpl) RefreshTagUpdateTimeByIDs(db *gorm.DB, ids []int64) error {
	return persistence.RefreshTagUpdateTimeByIDs(db, ids)
}

func (c ContentRepoImpl) SelectTagsByArticleID(db *gorm.DB, articleID int64) ([]*dto.Tag, error) {
	// 获取标签
	tagIDs, err := persistence.SelectTagIDsByArticleID(db, articleID)
	if err != nil {
		return nil, fmt.Errorf("select tag_id error:[%v]", err)
	}
	tagMap, err := persistence.MSelectTagByID(db, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("select tag error:[%v]", err)
	}
	return dto.NewTagsByModelMap(tagMap), nil
}

func (c ContentRepoImpl) UpdateArticleTagUpdateAtByArticleID(db *gorm.DB, id int64, publishAt *time.Time) error {
	return persistence.UpdateArticleTagUpdateAtByArticleID(db, id, publishAt)
}

func (c ContentRepoImpl) CreateTag(db *gorm.DB, tag *dto.Tag) (*dto.Tag, error) {
	created, err := persistence.CreateTag(db, &model.Tag{
		ID:      tag.ID,
		TagName: tag.TagName,
	})
	if err != nil {
		return nil, err
	}
	return dto.NewTagByModel(created), nil
}

func (c ContentRepoImpl) UpdateTag(db *gorm.DB, tag *dto.Tag) error {
	return persistence.UpdateTagByID(db, tag.ID, &model.Tag{
		TagName: tag.TagName,
	})
}

func (c ContentRepoImpl) DeleteTagByID(db *gorm.DB, id int64) error {
	return persistence.DeleteTagByID(db, id)
}

func (c ContentRepoImpl) DeleteArticleTagRelationByTagID(db *gorm.DB, tagID int64) error {
	return persistence.DeleteArticleTagRelationByTagID(db, tagID)
}

func (c ContentRepoImpl) GetTagListByPage(db *gorm.DB, keyword *string, page *int32, size *int32) ([]*dto.Tag, int64, error) {
	tags, err := persistence.GetTagListByPage(db, keyword, page, size)
	if err != nil {
		return nil, 0, err
	}
	tagsDTO := dto.NewTagsByModel(tags)
	tagIDs := make([]int64, len(tags))
	for i, tag := range tags {
		tagIDs[i] = tag.ID
	}
	articleCount, err := persistence.MGetTagArticleCountByTagIDs(db, tagIDs, false)
	if err != nil {
		return nil, 0, err
	}
	for _, tag := range tagsDTO {
		tag.Count = articleCount[tag.ID]
	}
	count, err := persistence.GetAllTagCount(db)
	if err != nil {
		return nil, 0, err
	}
	return tagsDTO, count, nil
}

func (c ContentRepoImpl) SelectArticleIDsByTagNames(db *gorm.DB, tagNames []string) ([]int64, error) {
	tagMap, err := persistence.MSelectTagByName(db, tagNames)
	if err != nil {
		return nil, err
	}
	var tagIDs []int64
	for _, category := range tagMap {
		tagIDs = append(tagIDs, category.ID)
	}
	if len(tagIDs) == 0 {
		return make([]int64, 0), nil
	}
	postIDs, err := persistence.SelectArticleIDsByTagIDs(db, tagIDs)
	if err != nil {
		return nil, err
	}
	return postIDs, nil
}
