package entity

import (
	"context"

	"github.com/spf13/cast"

	"my_blog/biz/domain/repo"
	"my_blog/biz/hertz_gen/blog/page"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/infra/render"
	"my_blog/biz/infra/session"
)

type ArticlePage struct {
	ID      int64
	Author  string
	Slug    string
	Article *Article
	Prev    *ArticleMeta
	Next    *ArticleMeta
}

func (a *ArticlePage) FetchContent(ctx context.Context) error {
	if a.Slug != "" {
		return a.fetchBySlug(ctx)
	}
	logger := log.GetLoggerWithCtx(ctx).WithField("article_id", a.ID)
	articleDTO, err := repo.GetContentRepo().Cache().GetArticle(ctx, a.ID)
	if err != nil {
		// 拿不到会返回 not found
		return err
	}
	a.Article = NewArticleByDTO(articleDTO, nil, nil)
	// 拿分类和标签
	a.Article.FetchCategoryAndTags(ctx)
	// 拿UV
	a.Article.FetchUV(ctx)
	// 拿上一篇和下一篇
	if articleDTO.PrevID != nil {
		prev, err := repo.GetContentRepo().Cache().GetArticleMeta(ctx, *articleDTO.PrevID)
		if err != nil {
			logger.Warnf("get prev article fail")
		}
		a.Prev = newArticleMetaByDTO(prev)
	}
	if articleDTO.NextID != nil {
		next, err := repo.GetContentRepo().Cache().GetArticleMeta(ctx, *articleDTO.NextID)
		if err != nil {
			logger.Warnf("get prev article fail")
		}
		a.Next = newArticleMetaByDTO(next)
	}
	a.Author = articleDTO.CreateUser.Nickname
	return nil
}

func (a *ArticlePage) fetchBySlug(ctx context.Context) error {
	articleDTO, err := repo.GetContentRepo().Cache().GetArticleBySlug(ctx, a.Slug)
	if err != nil {
		// 拿不到会返回 not found
		return err
	}
	article := &Article{}
	article.FillByDTO(articleDTO, nil, nil)
	a.Article = article
	a.Author = articleDTO.CreateUser.Nickname
	return nil
}

func (a *ArticlePage) PackPostPageResponse() *page.PostPageResponse {
	prevPage := &page.PostNav{}
	nextPage := &page.PostNav{}
	if a.Prev != nil {
		prevPage = &page.PostNav{
			ID:    cast.ToString(a.Prev.ID),
			Title: a.Prev.Title,
		}
	}
	if a.Next != nil {
		nextPage = &page.PostNav{
			ID:    cast.ToString(a.Next.ID),
			Title: a.Next.Title,
		}
	}
	categoryList := make([]*page.TermListItem, 0, len(a.Article.Categories))
	for _, category := range a.Article.Categories {
		categoryList = append(categoryList, &page.TermListItem{
			Name: category.CategoryName,
			Slug: category.Slug,
		})
	}
	return &page.PostPageResponse{
		Title: a.Article.Title,
		Info: &page.PostInfo{
			Author:       a.Author,
			PublishAt:    misc.GetPublishAtStr(a.Article.PublishAt),
			UV:           cast.ToString(a.Article.UV),
			WordCount:    misc.GetWordCount(a.Article.Content),
			CategoryList: categoryList,
		},
		Content:  render.RenderContent(a.Article.Content),
		Tags:     a.Article.Tags.ToStringList(),
		PrevPage: prevPage,
		NextPage: nextPage,
		Meta: resp.NewSuccessPageMeta(
			misc.GetPostPageTitle(a.Article.Title),
			misc.GetPostPageDescription(a.Article.Content),
			page.PageTypePost,
		),
	}
}

func (a *ArticlePage) PackPagePageResponse() *page.PagePageResponse {
	return &page.PagePageResponse{
		Title: a.Article.Title,
		Info: &page.PageInfo{
			Author:    a.Author,
			PublishAt: misc.GetPublishAtStr(a.Article.PublishAt),
			WordCount: misc.GetWordCount(a.Article.Content),
		},
		Content: render.RenderContent(a.Article.Content),
		Meta: resp.NewSuccessPageMeta(
			misc.GetPostPageTitle(a.Article.Title),
			misc.GetPostPageDescription(a.Article.Content),
			page.PageTypePage,
		),
	}
}

func (a *ArticlePage) View(ctx context.Context) {
	misc.SafeGo(ctx, func() {
		logger := log.GetLoggerWithCtx(ctx)
		vSession, err := session.NewVisitorSessionFromCtx(ctx)
		if err != nil {
			logger.Warnf("new visitor session error:[%v]", err)
			return
		}
		if vSession.CheckPostHasVisited(a.ID) {
			return
		}
		vSession.SetPostHasVisited(a.ID)
		err = repo.GetContentRepo().Cache().IncrPostUV(ctx, a.ID)
		if err != nil {
			logger.Warnf("ub incr error:[%v]", err)
			return
		}
		return
	})
	return
}
