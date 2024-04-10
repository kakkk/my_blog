package dto

import (
	"fmt"
	"time"

	"github.com/spf13/cast"

	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
	"my_blog/biz/model/blog/page"
)

type Article struct {
	ID           int64                `json:"id"`
	Title        string               `json:"title"`
	Content      string               `json:"content"`
	ArticleType  common.ArticleType   `json:"article_type"`
	Status       common.ArticleStatus `json:"status"`
	CreateUserID int64                `json:"create_user_id"`
	UV           int64                `json:"uv"`
	UpdateAt     time.Time            `json:"update_at"`
	PublishAt    *time.Time           `json:"publish_at"`
	Tags         []*Tag               `json:"tags"`
	Categories   []*Category          `json:"categories"`
	CreateUser   *User                `json:"create_user"`
	PrevID       *int64               `json:"prev_id"`
	NextID       *int64               `json:"next_id"`
}

type ArticleMeta struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Info        string `json:"info"`
	Description string `json:"description"`
	Abstract    string `json:"abstract"`
}

func (p *ArticleMeta) ToPostItem() *page.PostItem {
	return &page.PostItem{
		ID:       cast.ToString(p.ID),
		Title:    p.Title,
		Abstract: p.Abstract,
		Info:     p.Info,
	}
}

type ArticleMetas []*ArticleMeta

func (a ArticleMetas) ToPostItems() []*page.PostItem {
	result := make([]*page.PostItem, 0, len(a))
	for _, meta := range a {
		result = append(result, meta.ToPostItem())
	}
	return result
}

func (a ArticleMetas) PackPostListPageResp(currentPage int64, hasMore bool, pageType, name, slug string) *page.PostListPageResp {
	prev, next := "", ""
	if currentPage != 1 {
		prev = cast.ToString(currentPage - 1)
	}
	if hasMore {
		next = cast.ToString(currentPage + 1)
	}
	title, description := "", ""
	switch pageType {
	case page.PageTypeIndex:
		title = fmt.Sprintf("%v - %v", config.GetBlogName(), config.GetBlogSubTitle())
		description = config.GetBlogDescription()
	case page.PageTypePostList:
		title = fmt.Sprintf("第%v页 - %v - %v", currentPage, config.GetBlogName(), config.GetBlogSubTitle())
		description = config.GetBlogDescription()
	case page.PageTypeCategoryPostList:
		if currentPage == 1 {
			title = fmt.Sprintf("分类 %v 下的文章 - %v", name, config.GetBlogName())
		} else {
			title = fmt.Sprintf("分类 %v 下的文章 - 第%v页 - %v", name, currentPage, config.GetBlogName())
		}
		description = config.GetBlogDescription()
	case page.PageTypeTagPostList:
		if currentPage == 1 {
			title = fmt.Sprintf("标签 %v 下的文章 - %v", name, config.GetBlogName())
		} else {
			title = fmt.Sprintf("标签 %v 下的文章 - 第%v页 - %v", name, currentPage, config.GetBlogName())
		}
		description = config.GetBlogDescription()
	}
	return &page.PostListPageResp{
		Name:     name,
		Slug:     slug,
		PostList: a.ToPostItems(),
		PrevPage: prev,
		NextPage: next,
		Meta: &page.PageMeta{
			Title:       title,
			Description: description,
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			PageType:    pageType,
		},
	}
}

func NewArticleByModel(a *model.Article) *Article {
	return &Article{
		ID:           a.ID,
		Title:        a.Title,
		Content:      a.Content,
		ArticleType:  a.ArticleType,
		Status:       a.Status,
		CreateUserID: a.CreateUser,
		UV:           a.UV,
		UpdateAt:     a.UpdateAt,
		PublishAt:    a.PublishAt,
	}
}

func (a *Article) ToArticleMeta() *ArticleMeta {
	var editorName string
	// 降级
	if a.CreateUser == nil {
		editorName = config.GetDefaultUserName()
	} else {
		editorName = a.CreateUser.Nickname
	}
	return &ArticleMeta{
		ID:          a.ID,
		Title:       a.Title,
		Description: misc.GetPostPageDescription(a.Content),
		Info:        misc.GetPostInfo(editorName, *a.PublishAt, a.Content),
		Abstract:    misc.GetPostMetaAbstract(a.Content),
	}
}

type Articles []*Article

func (a Articles) ToRespList() []*api.PostListItem {
	result := make([]*api.PostListItem, 0, len(a))
	for _, article := range a {
		publishAt := int64(0)
		if article.PublishAt != nil {
			publishAt = article.PublishAt.Unix()
		}
		result = append(result, &api.PostListItem{
			ID:           article.ID,
			Title:        article.Title,
			CategoryList: Categories(article.Categories).ToStringList(),
			Editor:       article.CreateUser.Nickname,
			Status:       article.Status,
			UV:           article.UV,
			UpdateAt:     article.UpdateAt.Unix(),
			PublishAt:    publishAt,
		})
	}
	return result
}
