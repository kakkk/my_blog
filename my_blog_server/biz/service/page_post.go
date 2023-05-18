package service

import (
	"context"

	"github.com/spf13/cast"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/errorx"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/dto"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/storage"
)

func PostPage(ctx context.Context, req *page.PostPageRequest) (rsp *page.PostPageResponse, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx).WithField("post_id", req.GetID())
	rsp = page.NewPostPageResponse()
	defer utils.Recover(ctx, func() {
		pErr = errorx.NewInternalErrPageError()
		return
	})()
	// 获取post
	post, err := storage.GetArticleEntityStorage().Get(ctx, req.GetID())
	if err != nil {
		if err == consts.ErrRecordNotFound {
			logger.Warnf("post not found")
			return nil, errorx.NewNotFoundErrPageError()
		}
		logger.Errorf("select post error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}

	// 获取上下篇文章
	prev, next := getPrevNextPostMeta(ctx, req.GetID())

	// 获取作者
	user, err := storage.GetUserEntityStorage().Get(ctx, post.CreateUser)
	if err != nil {
		logger.Warnf("select user error:[%v]", err)
	}

	// 获取标签
	tags, err := storage.GetPostTagListStorage().Get(ctx, req.GetID())
	if err != nil {
		logger.Warnf("select tags error:[%v]", err)
	}

	// 获取分类
	categories, err := storage.GetPostCategoryListStorage().Get(ctx, req.GetID())
	if err != nil {
		logger.Warnf("select tags error:[%v]", err)
	}

	return packGetPostPageResp(post, prev, next, user, tags, categories), nil

}

func getPrevNextPostMeta(ctx context.Context, id int64) (*dto.PostMeta, *dto.PostMeta) {
	logger := log.GetLoggerWithCtx(ctx).WithField("post_id", id)
	prevNext, err := storage.GetPostPrevNextStorage().Get(ctx, id)
	if err != nil {
		logger.Warnf("get pre next page error:[%v]", err)
	}
	var ids []int64
	if prevNext.Prev != nil {
		ids = append(ids, *prevNext.Prev)
	}
	if prevNext.Next != nil {
		ids = append(ids, *prevNext.Next)
	}
	var prev *dto.PostMeta
	var next *dto.PostMeta
	metas := storage.GetPostMetaStorage().MGet(ctx, ids)
	if prevNext.Prev != nil {
		prev = metas[*prevNext.Prev]
	}
	if prevNext.Next != nil {
		next = metas[*prevNext.Next]
	}
	return prev, next
}

func packGetPostPageResp(post *entity.Article, prev *dto.PostMeta, next *dto.PostMeta, editor *entity.User, tags []string, categories []*entity.Category) *page.PostPageResponse {
	var author string
	prevPage := &page.PostNav{}
	nextPage := &page.PostNav{}
	if editor != nil {
		author = editor.Nickname
	}
	if prev != nil {
		prevPage = &page.PostNav{
			ID:    cast.ToString(prev.ID),
			Title: prev.Title,
		}
	}
	if next != nil {
		nextPage = &page.PostNav{
			ID:    cast.ToString(next.ID),
			Title: next.Title,
		}
	}
	categoryList := make([]*page.TermListItem, 0, len(categories))
	for _, category := range categories {
		categoryList = append(categoryList, &page.TermListItem{
			Name: category.CategoryName,
			Slug: category.Slug,
		})
	}
	rsp := &page.PostPageResponse{
		Title: post.Title,
		Info: &page.PostInfo{
			Author:       author,
			PublishAt:    utils.GetPublishAtStr(post.PublishAt),
			UV:           cast.ToString(post.UV),
			WordCount:    utils.GetWordCount(post.Content),
			CategoryList: categoryList,
		},
		Content:  utils.RenderContent(post.Content),
		Tags:     tags,
		PrevPage: prevPage,
		NextPage: nextPage,
		Meta: resp.NewSuccessPageMeta(
			utils.GetPostPageTitle(post.Title),
			utils.GetPostPageDescription(post.Content),
			page.PageTypePost,
		),
	}
	return rsp
}
