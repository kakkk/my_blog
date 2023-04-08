package mysql

import (
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/common"

	"gorm.io/gorm"
)

func CreateArticle(db *gorm.DB, article *entity.Article) (*entity.Article, error) {
	article.CreateAt = time.Now()
	article.UpdateAt = time.Now()
	err := db.Model(&entity.Article{}).Create(article).Error
	if err != nil {
		return nil, parseError(err)
	}
	return article, nil
}

func SelectArticleByID(db *gorm.DB, id int64) (*entity.Article, error) {
	post := &entity.Article{}
	err := db.Model(&entity.Article{}).
		Where("id = ?", id).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(post).Error
	if err != nil {
		return nil, parseError(err)
	}
	return post, nil
}

func UpdateArticleByID(db *gorm.DB, id int64, article *entity.Article) error {
	err := db.Model(&entity.Article{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"title":     article.Title,
			"content":   article.Content,
			"update_at": time.Now(),
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func UpdateArticleStatusByID(db *gorm.DB, id int64, status common.ArticleStatus) error {
	err := db.Model(&entity.Article{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"article_status": status,
			"update_at":      time.Now(),
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func UpdateArticlePublishAtByID(db *gorm.DB, id int64, publishAt *time.Time) error {
	err := db.Model(&entity.Article{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"publish_at": publishAt,
			"update_at":  time.Now(),
		}).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func SearchPostListByLimit(db *gorm.DB, keyword *string, ids []int64, page *int32, size *int32) ([]*entity.Article, error) {
	var result []*entity.Article
	tx := db.Model(&entity.Article{})
	if keyword != nil && *keyword != "" {
		tx.Where("title like ?", *keyword+"%")
	}
	tx.Where("delete_flag = ?", common.DeleteFlag_Exist)
	tx.Where("article_type = ?", common.ArticleType_Post)
	if ids != nil && len(ids) > 0 {
		tx.Where("id in (?)", ids)
	}

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
		Order("id desc").
		Find(&result).Error
	if err != nil {
		return nil, parseError(err)
	}
	if len(result) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return result, nil
}

func SelectSearchPostCount(db *gorm.DB, keyword *string, ids []int64) (int64, error) {
	count := int64(0)
	tx := db.Model(&entity.Article{})
	if keyword != nil && *keyword != "" {
		tx.Where("title like ?", *keyword+"%")
	}
	if ids != nil && len(ids) > 0 {
		tx.Where("id in (?)", ids)
	}
	err := tx.Where("delete_flag = ?", common.DeleteFlag_Exist).
		Where("article_type = ?", common.ArticleType_Post).
		Count(&count).Error
	if err != nil {
		return 0, parseError(err)
	}
	return count, nil
}

func DeleteArticleByID(db *gorm.DB, id int64) error {
	err := db.Model(&entity.Article{}).
		Where("id = ?", id).
		Updates(
			map[string]any{
				"article_status": common.ArticleStatus_DELETE,
				"delete_flag":    common.DeleteFlag_Delete,
				"update_at":      time.Now(),
			},
		).Error
	if err != nil {
		return parseError(err)
	}
	return nil
}

func SelectPostWithPublishByID(db *gorm.DB, id int64) (*entity.Article, error) {
	post := &entity.Article{}
	err := db.Model(&entity.Article{}).
		Where("id = ?", id).
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(post).Error
	if err != nil {
		return nil, parseError(err)
	}
	return post, nil
}

func SelectPrevNextPostByPublishAt(db *gorm.DB, publishAt time.Time) (*entity.Article, *entity.Article, error) {
	prev := &entity.Article{}
	next := &entity.Article{}
	prevErr := db.Model(&entity.Article{}).
		Where("publish_at > ?", publishAt).
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Order("publish_at asc").
		First(prev).Error
	nextErr := db.Model(&entity.Article{}).
		Where("publish_at < ?", publishAt).
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Order("publish_at desc").
		First(next).Error
	if prevErr != nil || nextErr != nil {
		return prev, next, fmt.Errorf("db err: prev:[%v], next:[%v]", parseError(prevErr), parseError(nextErr))
	}
	return prev, next, nil
}

func MSelectPostWithPublishByIDs(db *gorm.DB, ids []int64) (map[int64]*entity.Article, error) {
	var posts []*entity.Article
	err := db.Model(&entity.Article{}).
		Where("id in ?", ids).
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&posts).Error
	if err != nil {
		return nil, parseError(err)
	}
	if len(posts) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	result := make(map[int64]*entity.Article, len(posts))
	for _, post := range posts {
		result[post.ID] = post
	}
	return result, nil
}

func SelectPostOrderList(db *gorm.DB) ([]int64, error) {
	// 当前数据量直接拉取全量数据，后续可以加上limit分批次查询
	var order []int64
	err := db.Model(&entity.Article{}).
		Select("id").
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Order("publish_at desc").
		Find(&order).Error
	if err != nil {
		return nil, parseError(err)
	}
	if len(order) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}
