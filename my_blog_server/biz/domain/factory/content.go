package factory

import (
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
)

type ContentFactory struct{}

var contentFactory = ContentFactory{}

func GetContentFactory() ContentFactory {
	return contentFactory
}

func (ContentFactory) NewArticleByID(id int64) *entity.Article {
	return &entity.Article{
		ID: id,
	}
}

func (ContentFactory) NewArticleByDTO(article *dto.Article) *entity.Article {
	return &entity.Article{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		ArticleType: article.ArticleType,
		Status:      article.Status,
		Slug:        article.Slug,
		CreateUser:  article.CreateUserID,
		UV:          article.UV,
		PublishAt:   article.PublishAt,
	}
}
