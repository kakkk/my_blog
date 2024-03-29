package dto

import (
	"my_blog/biz/common/config"
	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/page"

	"github.com/spf13/cast"
)

type PostMeta struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Info        string `json:"info"`
	Description string `json:"description"`
	Abstract    string `json:"abstract"`
}

func (p *PostMeta) ToPostItem() *page.PostItem {
	return &page.PostItem{
		ID:       cast.ToString(p.ID),
		Title:    p.Title,
		Abstract: p.Abstract,
		Info:     p.Info,
	}
}

func NewArticleMetaByEntity(post *entity.Article, editor *entity.User) *PostMeta {
	var editorName string
	// 降级
	if editor == nil {
		editorName = config.GetDefaultUserName()
	} else {
		editorName = editor.Nickname
	}
	return &PostMeta{
		ID:          post.ID,
		Title:       post.Title,
		Description: utils.GetPostPageDescription(post.Content),
		Info:        utils.GetPostInfo(editorName, *post.PublishAt, post.Content),
		Abstract:    utils.GetPostMetaAbstract(post.Content),
	}
}
