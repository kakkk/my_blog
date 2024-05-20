package entity

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"
)

type Article struct {
	ID          int64
	Title       string
	Content     string
	ArticleType common.ArticleType
	Status      common.ArticleStatus
	Slug        string
	CreateUser  int64
	UV          int64
	PublishAt   *time.Time
	Tags        Tags
	Categories  Categories
}

func (a *Article) Create(tx *gorm.DB) error {
	createdArticle, err := repo.GetContentRepo().CreateArticle(tx, &dto.Article{
		Title:        a.Title,
		Content:      a.Content,
		ArticleType:  a.ArticleType,
		Status:       a.Status,
		Slug:         a.Slug,
		CreateUserID: a.CreateUser,
		PublishAt:    a.PublishAt,
	})
	if err != nil {
		return fmt.Errorf("create article error:[%v]", err)
	}
	a.ID = createdArticle.ID
	return nil
}

func (a *Article) Update(tx *gorm.DB) error {
	err := repo.GetContentRepo().UpdateArticleByID(tx, a.ID, &dto.Article{
		Title:   a.Title,
		Content: a.Content,
		Slug:    a.Slug,
	})
	if err != nil {
		return fmt.Errorf("create article error:[%v]", err)
	}
	return nil
}

func (a *Article) Get(db *gorm.DB) error {
	// 获取post
	article, err := repo.GetContentRepo().SelectArticleByID(db, a.ID)
	if err != nil {
		return err
	}
	a.FillByDTO(article, nil, nil)
	return nil
}

func (a *Article) FillByDTO(article *dto.Article, categories []*dto.Category, tags []*dto.Tag) {
	if article != nil {
		a.ID = article.ID
		a.Title = article.Title
		a.Content = article.Content
		a.ArticleType = article.ArticleType
		a.Status = article.Status
		a.Slug = article.Slug
		a.CreateUser = article.CreateUserID
		a.UV = article.UV
		a.PublishAt = article.PublishAt
	}
}

func (a *Article) LinkCategoriesAndTags(tx *gorm.DB) error {
	// 解除所有关联
	err := repo.GetContentRepo().DeleteArticleCategoryRelationByArticleID(tx, a.ID)
	if err != nil {
		return fmt.Errorf("delete article category relation error:[%v]", err)
	}

	// 添加分类
	var categoryRelation []*model.ArticleCategory
	for _, category := range a.Categories {
		categoryRelation = append(categoryRelation, &model.ArticleCategory{
			PostID:     a.ID,
			CategoryID: category.ID,
			PublishAt:  a.PublishAt,
		})
	}
	err = repo.GetContentRepo().UpsertArticleCategoryRelation(tx, categoryRelation)
	if err != nil {
		return fmt.Errorf("upsert article_category error:[%v]", err)
	}

	// 先清除标签与文章关联关系
	err = repo.GetContentRepo().DeleteArticleTagRelationByArticleID(tx, a.ID)
	if err != nil {
		return fmt.Errorf("delete tag error:[%v]", err)
	}

	// 文章添加标签
	var tagRelation []*model.ArticleTag
	for _, tag := range a.Tags {
		tagRelation = append(tagRelation, &model.ArticleTag{
			PostID:    a.ID,
			TagID:     tag.ID,
			PublishAt: a.PublishAt,
		})
	}
	err = persistence.UpsertArticleTagRelation(tx, tagRelation)
	if err != nil {
		return fmt.Errorf("upsert article_tag error:[%v]", err)
	}

	// 刷新标签状态
	err = persistence.RefreshTagUpdateTimeByIDs(tx, a.Tags.IDs())
	if err != nil {
		return fmt.Errorf("refresh tag update_time error:[%v]", err)
	}

	return nil
}

func (a *Article) UpdateStatus(db *gorm.DB) error {
	article, err := repo.GetContentRepo().SelectArticleByID(db, a.ID)
	if err != nil {
		return fmt.Errorf("select article by id error:[%w]", err)
	}

	// 新发布文章
	if a.Status == common.ArticleStatus_PUBLISH && article.PublishAt == nil {
		publishAt := time.Now()
		// 更新文章发布时间
		err = repo.GetContentRepo().UpdateArticlePublishAtByID(db, a.ID, &publishAt)
		if err != nil {
			return fmt.Errorf("update article publish_at by id error:[%w]", err)
		}
		// 更新分类发布时间
		err = repo.GetContentRepo().UpdateArticleCategoryUpdateAtByArticleID(db, a.ID, &publishAt)
		if err != nil {
			return fmt.Errorf("update article_category publish_at by id error:[%w]", err)
		}
		// 更新标签发布时间
		err = repo.GetContentRepo().UpdateArticleTagUpdateAtByArticleID(db, a.ID, &publishAt)
		if err != nil {
			return fmt.Errorf("update article_tag publish_at by id error:[%w]", err)
		}
	}

	// 更新状态
	err = repo.GetContentRepo().UpdateArticleStatusByID(db, a.ID, a.Status)
	if err != nil {
		return fmt.Errorf("update article status by id error:[%w]", err)
	}

	return nil

}

func (a *Article) Delete(ctx context.Context) error {
	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		err := repo.GetContentRepo().DeleteArticleByID(tx, a.ID)
		if err != nil {
			return fmt.Errorf("delete article error:[%v]", err)
		}
		err = repo.GetContentRepo().DeleteArticleCategoryRelationByArticleID(tx, a.ID)
		if err != nil {
			return fmt.Errorf("delete article_category relation error:[%v]", err)
		}
		err = repo.GetContentRepo().DeleteArticleTagRelationByArticleID(tx, a.ID)
		if err != nil {
			return fmt.Errorf("delete article_tag relation error:[%v]", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *Article) FetchUV(ctx context.Context) {
	uv, err := repo.GetContentRepo().Cache().GetPostUV(ctx, a.ID)
	if err != nil {
		// 获取不到最新的UV，不影响正常显示，打日志即可
		log.GetLoggerWithCtx(ctx).
			WithField("article_id", a.ID).
			Warnf("fetch uv fail: %v", err)
		return
	}
	a.UV = uv
}

func (a *Article) FetchCategoryAndTags(ctx context.Context) {
	logger := log.GetLoggerWithCtx(ctx).WithField("article_id", a.ID)
	categories, err := repo.GetContentRepo().Cache().GetCategoriesByArticleID(ctx, a.ID)
	if err != nil {
		logger.Warnf("fetch categories fail: %v", err)
	}
	tags, err := repo.GetContentRepo().Cache().GetTagsByArticleID(ctx, a.ID)
	if err != nil {
		logger.Warnf("fetch tags fail: %v", err)
	}
	a.Categories = NewCategoriesByDTO(categories)
	a.Tags = NewTagsByDTO(tags)
}

func NewArticleByDTO(article *dto.Article, categories []*dto.Category, tags []*dto.Tag) *Article {
	return &Article{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		ArticleType: article.ArticleType,
		Status:      article.Status,
		Slug:        article.Slug,
		CreateUser:  article.CreateUserID,
		UV:          article.UV,
		PublishAt:   article.PublishAt,
		Tags:        NewTagsByDTO(tags),
		Categories:  NewCategoriesByDTO(categories),
	}
}

type ArticleMeta struct {
	ID          int64
	Title       string
	Info        string
	Description string
	Abstract    string
}

func newArticleMetaByDTO(am *dto.ArticleMeta) *ArticleMeta {
	return &ArticleMeta{
		ID:          am.ID,
		Title:       am.Title,
		Info:        am.Info,
		Description: am.Description,
		Abstract:    am.Abstract,
	}
}

type Articles struct {
	Articles []*Article
}

func (a *Articles) MGetMyIDs(ctx context.Context, ids []int64) error {
	articles, err := repo.GetContentRepo().SelectArticleByIDs(mysql.GetDB(ctx), ids)
	if err != nil {
		return err
	}
	for _, article := range articles {
		a.Articles = append(a.Articles, NewArticleByDTO(article, nil, nil))
	}
	return nil
}

func (a *Articles) GetID2Title() map[int64]string {
	result := make(map[int64]string)
	for _, article := range a.Articles {
		result[article.ID] = article.Title
	}
	return result
}
