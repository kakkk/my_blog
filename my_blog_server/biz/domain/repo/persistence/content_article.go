package persistence

import (
	"time"

	"my_blog/biz/consts"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"

	"gorm.io/gorm"
)

func CreateArticle(db *gorm.DB, article *model.Article) (*model.Article, error) {
	article.CreateAt = time.Now()
	article.UpdateAt = time.Now()
	err := db.Model(&model.Article{}).Create(article).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return article, nil
}

func SelectArticleByID(db *gorm.DB, id int64) (*model.Article, error) {
	post := &model.Article{}
	err := db.Model(&model.Article{}).
		Where("id = ?", id).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(post).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return post, nil
}

func SelectArticleBySlug(db *gorm.DB, slug string) (*model.Article, error) {
	article := &model.Article{}
	err := db.Model(&model.Article{}).
		Where("slug = ?", slug).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		First(article).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return article, nil
}

func UpdateArticleByID(db *gorm.DB, id int64, article *model.Article) error {
	err := db.Model(&model.Article{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"title":     article.Title,
			"content":   article.Content,
			"slug":      article.Slug,
			"update_at": time.Now(),
		}).Error
	if err != nil {
		return mysql.ParseError(err)
	}
	return nil
}

func UpdateArticleStatusByID(db *gorm.DB, id int64, status common.ArticleStatus) error {
	err := db.Model(&model.Article{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"article_status": status,
			"update_at":      time.Now(),
		}).Error
	if err != nil {
		return mysql.ParseError(err)
	}
	return nil
}

func UpdateArticlePublishAtByID(db *gorm.DB, id int64, publishAt *time.Time) error {
	err := db.Model(&model.Article{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"publish_at": publishAt,
			"update_at":  time.Now(),
		}).Error
	if err != nil {
		return mysql.ParseError(err)
	}
	return nil
}

func SearchPostListByLimit(db *gorm.DB, keyword *string, ids []int64, page *int32, size *int32) ([]*model.Article, error) {
	var result []*model.Article
	tx := db.Model(&model.Article{})
	if keyword != nil && *keyword != "" {
		tx.Where("title like ?", *keyword+"%")
	}
	tx.Where("delete_flag = ?", common.DeleteFlag_Exist)
	tx.Where("article_type = ?", common.ArticleType_Post)
	if ids != nil && len(ids) > 0 {
		tx.Where("id in (?)", ids)
	}

	offset, limit := 0, mysql.DefaultPageLimit
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
		return nil, mysql.ParseError(err)
	}
	if len(result) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return result, nil
}

func SelectSearchPostCount(db *gorm.DB, keyword *string, ids []int64) (int64, error) {
	count := int64(0)
	tx := db.Model(&model.Article{})
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
		return 0, mysql.ParseError(err)
	}
	return count, nil
}

func DeleteArticleByID(db *gorm.DB, id int64) error {
	err := db.Model(&model.Article{}).
		Where("id = ?", id).
		Updates(
			map[string]any{
				"article_status": common.ArticleStatus_DELETE,
				"delete_flag":    common.DeleteFlag_Delete,
				"update_at":      time.Now(),
			},
		).Error
	if err != nil {
		return mysql.ParseError(err)
	}
	return nil
}

func MSelectPostWithPublishByIDs(db *gorm.DB, ids []int64) (map[int64]*model.Article, error) {
	var posts []*model.Article
	err := db.Model(&model.Article{}).
		Where("id in ?", ids).
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&posts).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	if len(posts) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	result := make(map[int64]*model.Article, len(posts))
	for _, post := range posts {
		result[post.ID] = post
	}
	return result, nil
}

func SelectPostOrderList(db *gorm.DB) ([]int64, error) {
	// 当前数据量直接拉取全量数据，后续可以加上limit分批次查询
	var order []int64
	err := db.Model(&model.Article{}).
		Select("id").
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Order("publish_at desc").
		Find(&order).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	if len(order) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return order, nil
}

func SelectPostIDsByCategoryID(db *gorm.DB, cID int64) ([]int64, error) {
	var list []int64
	err := db.Model(&model.ArticleCategory{}).
		Select("article_id").
		Where("category_id = ?", cID).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Where("publish_at is not null").
		Order("publish_at desc").
		Find(&list).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	if len(list) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return list, nil
}

func SelectAllPages(db *gorm.DB) ([]*model.Article, error) {
	var articles []*model.Article
	err := db.Model(&model.Article{}).
		Where("article_type = ?", common.ArticleType_Page).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		Find(&articles).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	return articles, nil
}
