package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"my_blog/biz/consts"
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/common"
)

type ContentService struct{}

func GetContentService() *ContentService {
	return &ContentService{}
}

// CreateArticle 创建文章
func (svc *ContentService) CreateArticle(ctx context.Context, articleDTO *dto.Article, categoryIDs []int64, tagNames []string) (*dto.Article, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(
		logrus.Fields{
			"title": articleDTO.Title,
		},
	)
	var createdArticle *dto.Article
	// 开启事务
	txErr := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 先查询一下分类，不存在会返回NotFound
		categoriesDTO, err := repo.GetContentRepo().SelectCategoryByIDs(tx, categoryIDs)
		if err != nil {
			return err
		}

		// 创建不存在的标签，恢复已删除的标签，获取已存在的标签
		tagsDTO, err := repo.GetContentRepo().GetOrCreateTagsByNames(tx, tagNames)
		if err != nil {
			return err
		}

		article := entity.NewArticleByDTO(articleDTO, categoriesDTO, tagsDTO)

		// 创建文章
		err = article.Create(tx)
		if err != nil {
			return err
		}

		// 创建分类、标签和文章关联关系
		err = article.LinkCategoriesAndTags(tx)
		if err != nil {
			return err
		}

		// 提交事务
		return nil
	})

	if txErr != nil {
		logger.Errorf("create article with transaction fail:[%v]", txErr)
		return nil, fmt.Errorf("create article with transaction fail:[%v]", txErr)
	}

	return createdArticle, nil
}

// UpdateArticle 更新文章
func (svc *ContentService) UpdateArticle(ctx context.Context, articleDTO *dto.Article, categoryIDs []int64, tagNames []string) error {
	logger := log.GetLoggerWithCtx(ctx).WithFields(
		logrus.Fields{
			"title": articleDTO.Title,
		},
	)
	// 开启事务
	txErr := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 先查询一下分类，不存在会返回NotFound
		categoriesDTO, err := repo.GetContentRepo().SelectCategoryByIDs(tx, categoryIDs)
		if err != nil {
			return err
		}

		// 创建不存在的标签，恢复已删除的标签，获取已存在的标签
		tagsDTO, err := repo.GetContentRepo().GetOrCreateTagsByNames(tx, tagNames)
		if err != nil {
			return err
		}

		articleDTO, err := repo.GetContentRepo().SelectArticleByID(tx, articleDTO.ID)
		if err != nil {
			return err
		}
		article := entity.NewArticleByDTO(articleDTO, categoriesDTO, tagsDTO)

		// 更新文章
		err = article.Update(tx)
		if err != nil {
			return err
		}

		// 创建分类、标签和文章关联关系
		err = article.LinkCategoriesAndTags(tx)
		if err != nil {
			return err
		}

		// 提交事务
		return nil
	})

	if txErr != nil {
		logger.Errorf("update articleDTO with transaction fail:[%v]", txErr)
		return fmt.Errorf("update articleDTO with transaction fail:[%v]", txErr)
	}
	return nil
}

// GetArticleFromDB 从DB查询文章
func (svc *ContentService) GetArticleFromDB(ctx context.Context, id int64) (*dto.Article, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": id,
	})
	db := mysql.GetDB(ctx)

	// 获取post
	article, err := repo.GetContentRepo().SelectArticleByID(db, id)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			logger.Warnf("post not found, post_id:[%v]", id)
			return nil, err
		}
		logger.Errorf("select post error:[%v]", err)
		return nil, fmt.Errorf("select post error:[%v]", err)
	}

	// 获取分类
	categories, err := repo.GetContentRepo().GetCategoriesByArticleID(db, id)
	if err != nil {
		logger.Errorf("get categories fail: %v", err)
		return nil, fmt.Errorf("get categories fail: %v", err)
	}
	article.Categories = categories

	// 标签
	tags, err := repo.GetContentRepo().SelectTagsByArticleID(db, id)
	if err != nil {
		logger.Errorf("select tag list error:[%v]", err)
		return nil, fmt.Errorf("select tag list error:[%v]", err)
	}
	article.Tags = tags

	return article, nil
}

// UpdateArticleStatus 更新文章状态
func (svc *ContentService) UpdateArticleStatus(ctx context.Context, id int64, status common.ArticleStatus) error {
	article := &entity.Article{
		ID:     id,
		Status: status,
	}
	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		return article.UpdateStatus(tx)
	})
	if err != nil {
		return fmt.Errorf("update article status fail:[%w]", err)
	}
	return nil
}

// SearchArticleList 搜索文章列表
func (svc *ContentService) SearchArticleList(ctx context.Context, keyword *string, categories []string, tags []string, page *int32, size *int32) ([]*dto.Article, *dto.Pagination, error) {
	var searchIDs []int64
	searchByID := false
	db := mysql.GetDB(ctx)
	// 搜索分类下所有文章
	if len(categories) > 0 {
		searchByID = true
		ids, err := repo.GetContentRepo().SelectArticleIDsByCategoryNames(db, categories)
		if err != nil {
			return nil, nil, err
		}
		searchIDs = append(searchIDs, ids...)
	}
	// 搜索标签下所有文章
	if len(tags) > 0 {
		searchByID = true
		ids, err := repo.GetContentRepo().SelectArticleIDsByTagNames(db, tags)
		if err != nil {
			return nil, nil, err
		}
		if len(searchIDs) > 0 {
			// 取交集
			searchIDs = misc.IntersectInt64Slice(searchIDs, ids)
		} else {
			searchIDs = append(searchIDs, ids...)
		}
	}

	if searchByID && len(searchIDs) == 0 {
		return make([]*dto.Article, 0), dto.NewPagination(0), nil
	}

	// 获取总数
	total, err := repo.GetContentRepo().GetSearchPostTotal(db, keyword, searchIDs)
	if err != nil {
		return nil, nil, err
	}

	// 搜索文章
	articles, err := repo.GetContentRepo().SearchPostListByLimit(db, keyword, searchIDs, page, size)
	if err != nil {
		return nil, nil, err
	}

	return articles, dto.NewPagination(total), nil
}
