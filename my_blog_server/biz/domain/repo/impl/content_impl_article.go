package impl

import (
	"time"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/repository/model"
)

func (c ContentRepoImpl) CreateArticle(db *gorm.DB, article *dto.Article) (*dto.Article, error) {
	createdArticle, err := persistence.CreateArticle(db, &model.Article{
		Title:       article.Title,
		Content:     article.Content,
		ArticleType: article.ArticleType,
		Status:      article.Status,
		Slug:        article.Slug,
		CreateUser:  article.CreateUserID,
		PublishAt:   article.PublishAt,
	})
	if err != nil {
		return nil, err
	}
	return dto.NewArticleByModel(createdArticle), nil
}

func (c ContentRepoImpl) UpdateArticleByID(db *gorm.DB, id int64, article *dto.Article) error {
	return persistence.UpdateArticleByID(db, id, &model.Article{
		Title:   article.Title,
		Content: article.Content,
		Slug:    article.Slug,
	})
}

func (c ContentRepoImpl) SelectArticleByID(db *gorm.DB, id int64) (*dto.Article, error) {
	article, err := persistence.SelectArticleByID(db, id)
	if err != nil {
		return nil, err
	}
	return dto.NewArticleByModel(article), nil
}

func (c ContentRepoImpl) UpdateArticlePublishAtByID(db *gorm.DB, id int64, publishAt *time.Time) error {
	return persistence.UpdateArticlePublishAtByID(db, id, publishAt)
}

func (c ContentRepoImpl) UpdateArticleStatusByID(db *gorm.DB, id int64, status common.ArticleStatus) error {
	return persistence.UpdateArticleStatusByID(db, id, status)
}

func (c ContentRepoImpl) DeleteArticleByID(db *gorm.DB, id int64) error {
	return persistence.DeleteArticleByID(db, id)
}

func (c ContentRepoImpl) GetSearchPostTotal(db *gorm.DB, keyword *string, ids []int64) (int64, error) {
	return persistence.SelectSearchPostCount(db, keyword, ids)
}

func (c ContentRepoImpl) SearchPostListByLimit(db *gorm.DB, keyword *string, ids []int64, page *int32, size *int32) ([]*dto.Article, error) {
	// 拿文章
	articles, err := persistence.SearchPostListByLimit(db, keyword, ids, page, size)
	if err != nil {
		return nil, err
	}
	articleIDs := make([]int64, len(articles))
	userIDSet := make(map[int64]struct{})
	for i, article := range articles {
		articleIDs[i] = article.ID
		userIDSet[article.CreateUser] = struct{}{}
	}
	// 拿用户
	users, err := persistence.MSelectUserByIDs(db, misc.MapKeys(userIDSet))
	if err != nil {
		return nil, err
	}
	// 拿分类
	articleIDToCategoryIDs, err := persistence.MSelectCategoryIDsByArticleIDs(db, articleIDs)
	if err != nil {
		return nil, err
	}
	categoryIDSet := make(map[int64]struct{})
	for _, cIDs := range articleIDToCategoryIDs {
		for _, id := range cIDs {
			categoryIDSet[id] = struct{}{}
		}
	}
	categories, err := persistence.MSelectCategoryByIDs(db, misc.MapKeys(categoryIDSet))
	if err != nil {
		return nil, err
	}
	// 拿标签
	articleIDToTagIDs, err := persistence.MSelectTagIDsByArticleID(db, articleIDs)
	if err != nil {
		return nil, err
	}
	tagIDSet := make(map[int64]struct{})
	for _, tIDs := range articleIDToTagIDs {
		for _, id := range tIDs {
			tagIDSet[id] = struct{}{}
		}
	}
	tags, err := persistence.MSelectTagByID(db, misc.MapKeys(tagIDSet))
	if err != nil {
		return nil, err
	}

	result := make([]*dto.Article, 0, len(articles))
	for _, article := range articles {
		articleDTO := dto.NewArticleByModel(article)
		// 分类
		categoryIDs := articleIDToCategoryIDs[article.ID]
		categoriesDTO := make([]*dto.Category, 0, len(categoryIDs))
		for _, cID := range categoryIDs {
			categoriesDTO = append(categoriesDTO, dto.NewCategoryByModel(categories[cID]))
		}
		articleDTO.Categories = categoriesDTO
		// 标签
		tagIDs := articleIDToTagIDs[article.ID]
		tagsDTO := make([]*dto.Tag, 0, len(tagIDs))
		for _, tID := range tagIDs {
			tagsDTO = append(tagsDTO, dto.NewTagByModel(tags[tID]))
		}
		articleDTO.Tags = tagsDTO
		articleDTO.CreateUser = dto.NewUserByModel(users[article.CreateUser])
		result = append(result, articleDTO)
	}
	return result, nil
}

func (c ContentRepoImpl) SelectAllPages(db *gorm.DB) ([]*dto.Article, error) {
	articles, err := persistence.SelectAllPages(db)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.Article, 0, len(articles))
	for _, article := range articles {
		result = append(result, dto.NewArticleByModel(article))
	}
	return result, nil
}
