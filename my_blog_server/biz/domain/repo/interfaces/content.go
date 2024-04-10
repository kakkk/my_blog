package interfaces

import (
	"context"
	"time"

	"gorm.io/gorm"

	"my_blog/biz/domain/dto"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/model/blog/common"
)

type ContentRepo interface {
	CreateCategory(db *gorm.DB, category *dto.Category) (*dto.Category, error)
	SelectCategoryOrder(db *gorm.DB) ([]int64, error)
	UpdateCategoryOrder(db *gorm.DB, order []int64) error
	GetCategoryList(db *gorm.DB, publish bool) ([]*dto.Category, error)
	SelectCategoryByIDs(db *gorm.DB, ids []int64) ([]*dto.Category, error)
	DeleteArticleCategoryRelationByArticleID(db *gorm.DB, articleID int64) error
	UpsertArticleCategoryRelation(db *gorm.DB, articleCategories []*model.ArticleCategory) error
	GetCategoriesByArticleID(db *gorm.DB, articleID int64) ([]*dto.Category, error)
	UpdateArticleCategoryUpdateAtByArticleID(db *gorm.DB, id int64, publishAt *time.Time) error
	SelectArticleIDsByCategoryNames(db *gorm.DB, categoryNames []string) ([]int64, error)

	SelectArticleByID(db *gorm.DB, id int64) (*dto.Article, error)
	CreateArticle(db *gorm.DB, article *dto.Article) (*dto.Article, error)
	UpdateArticleByID(db *gorm.DB, id int64, article *dto.Article) error
	UpdateArticlePublishAtByID(db *gorm.DB, id int64, publishAt *time.Time) error
	UpdateArticleStatusByID(db *gorm.DB, id int64, status common.ArticleStatus) error
	DeleteArticleByID(db *gorm.DB, id int64) error
	GetSearchPostTotal(db *gorm.DB, keyword *string, ids []int64) (int64, error)
	SearchPostListByLimit(db *gorm.DB, keyword *string, ids []int64, page *int32, size *int32) ([]*dto.Article, error)

	GetOrCreateTagsByNames(db *gorm.DB, names []string) ([]*dto.Tag, error)
	DeleteArticleTagRelationByArticleID(db *gorm.DB, articleID int64) error
	UpsertArticleTagRelation(db *gorm.DB, articleTags []*model.ArticleTag) error
	RefreshTagUpdateTimeByIDs(db *gorm.DB, ids []int64) error
	SelectTagsByArticleID(db *gorm.DB, articleID int64) ([]*dto.Tag, error)
	UpdateArticleTagUpdateAtByArticleID(db *gorm.DB, id int64, publishAt *time.Time) error
	CreateTag(db *gorm.DB, tag *dto.Tag) (*dto.Tag, error)
	UpdateTag(db *gorm.DB, tag *dto.Tag) error
	DeleteTagByID(db *gorm.DB, id int64) error
	DeleteArticleTagRelationByTagID(db *gorm.DB, tagID int64) error
	GetTagListByPage(db *gorm.DB, keyword *string, page *int32, size *int32) ([]*dto.Tag, int64, error)
	SelectArticleIDsByTagNames(db *gorm.DB, tagNames []string) ([]int64, error)

	Cache() ContentCache
}

type ContentCache interface {
	GetDefaultCategoryID() int64
	GetCategoryList(ctx context.Context) ([]*dto.Category, error)
	GetCategoriesByArticleID(ctx context.Context, id int64) ([]*dto.Category, error)

	GetArticle(ctx context.Context, id int64) (*dto.Article, error)
	GetArticleMeta(ctx context.Context, id int64) (*dto.ArticleMeta, error)
	GetArticlePostIDs(ctx context.Context) ([]int64, error)
	MGetArticleMeta(ctx context.Context, ids []int64) map[int64]*dto.ArticleMeta
	GetPostUV(ctx context.Context, id int64) (int64, error)
	IncrPostUV(ctx context.Context, id int64) error

	GetTagList(ctx context.Context) ([]*dto.Tag, error)
	GetTagsByArticleID(ctx context.Context, id int64) ([]*dto.Tag, error)
	GetArticleIDsByTagName(ctx context.Context, name string) ([]int64, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*dto.Category, error)
	GetArticleIDsByCategoryID(ctx context.Context, id int64) ([]int64, error)
}
