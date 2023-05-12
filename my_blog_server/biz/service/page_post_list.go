package service

import (
	"context"
	"fmt"

	"github.com/spf13/cast"

	"my_blog/biz/common/config"
	"my_blog/biz/common/errorx"
	"my_blog/biz/common/log"
	"my_blog/biz/common/utils"
	"my_blog/biz/dto"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/storage"
)

func PostListByPage(ctx context.Context, req *page.PostListPageRequest) (rsp *page.PostListPageResp, pErr *errorx.PageError) {
	logger := log.GetLoggerWithCtx(ctx)
	defer utils.Recover(ctx, func() {
		pErr = errorx.NewInternalErrPageError()
		return
	})()
	postOrderList, err := storage.GetPostOrderListStorage().Get(ctx)
	if err != nil {
		logger.Errorf("get post order list error:[%v]", err)
		return nil, errorx.NewFailErrPageError()
	}
	postIDs, hasMore := getIDsByPage(postOrderList, req.GetPage())
	if len(postIDs) == 0 {
		return nil, errorx.NewNotFoundErrPageError()
	}
	postMetas := storage.GetPostMetaStorage().MGet(ctx, postIDs)
	if len(postMetas) == 0 {
		return nil, errorx.NewNotFoundErrPageError()
	}
	return packPostListPageResp(req.GetPage(), hasMore, req.GetPageType(), "", "", utils.MapToList(postIDs, postMetas)), nil
}

func packPostListPageResp(currentPage int64, hasMore bool, pageType, name, slug string, metas []*dto.PostMeta) *page.PostListPageResp {
	itemList := make([]*page.PostItem, 0, len(metas))
	for _, meta := range metas {
		itemList = append(itemList, meta.ToPostItem())
	}
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
	}
	return &page.PostListPageResp{
		Name:     name,
		Slug:     slug,
		PostList: itemList,
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

func getIDsByPage(ids []int64, page int64) ([]int64, bool) {
	if page < 1 {
		return []int64{}, false
	}
	if page == 1 {
		if len(ids) < config.GetPageListSize() {
			return ids, false
		}
		hasMore := true
		if len(ids) <= config.GetPageListSize() {
			hasMore = false
		}
		ids = ids[0:config.GetPageListSize()]
		result := make([]int64, len(ids))
		copy(result, ids)

		return result, hasMore
	}
	begin := (page - 1) * config.GetPageListSizeI64()
	end := begin + config.GetPageListSizeI64()
	hasMore := true
	if len(ids) <= int(end) {
		hasMore = false
	}
	if len(ids) < int(begin) {
		return []int64{}, false
	}
	if len(ids) < int(end) {
		ids = ids[begin:]
		result := make([]int64, len(ids))
		copy(result, ids)
		return result, false
	}
	ids = ids[begin:end]
	result := make([]int64, len(ids))
	copy(result, ids)
	return result, hasMore
}
