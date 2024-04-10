package persistence

import (
	"gorm.io/gorm"

	"my_blog/biz/consts"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/common"
)

func SelectAllPublishedPostWithBatch(db *gorm.DB) ([]*model.Article, error) {
	var list []*model.Article
	var results []*model.Article
	err := db.Model(&model.Article{}).
		Where("article_type = ?", common.ArticleType_Post).
		Where("article_status = ?", common.ArticleStatus_PUBLISH).
		Where("delete_flag = ?", common.DeleteFlag_Exist).
		FindInBatches(&results, 5, func(tx *gorm.DB, batch int) error {
			for _, result := range results {
				list = append(list, result)
			}
			return nil
		}).Error
	if err != nil {
		return nil, mysql.ParseError(err)
	}
	if len(list) == 0 {
		return nil, consts.ErrRecordNotFound
	}
	return list, nil
}
