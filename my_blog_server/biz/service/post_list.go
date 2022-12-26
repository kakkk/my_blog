package service

import (
	"context"
	"errors"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/repository/mysql"
)

func GetPostListAPI(ctx context.Context, req *api.GetPostListAPIRequest) (rsp *api.GetPostListAPIResponse) {
	logger := log.GetLoggerWithCtx(ctx)
	rsp = &api.GetPostListAPIResponse{}
	db := mysql.GetDB(ctx)
	var posts []*entity.Article
	total := int64(0)
	searchByID := false // 是否通过ID搜索
	var searchIDs []int64
	// 分类
	if len(req.GetCategories()) > 0 {
		searchByID = true
		categoryMap, err := mysql.MSelectCategoryByNames(db, req.GetCategories())
		if err != nil {
			logger.Errorf("select category by name error:[%v]", err)
			rsp.BaseResp = resp.NewFailBaseResp()
			return
		}
		var categoryIDs []int64
		for _, category := range categoryMap {
			categoryIDs = append(categoryIDs, category.ID)
		}
		postIDs, err := mysql.SelectArticleIDsByCategoryIDs(db, categoryIDs)
		if err != nil {
			logger.Errorf("select post_id by category_id error:[%v]", err)
			rsp.BaseResp = resp.NewFailBaseResp()
			return
		}
		searchIDs = append(searchIDs, postIDs...)
	}
	// 标签
	if len(req.GetTags()) > 0 {
		searchByID = true
		tagMap, err := mysql.MSelectTagByName(db, req.GetTags())
		if err != nil {
			logger.Errorf("select tag by name error:[%v]", err)
			rsp.BaseResp = resp.NewFailBaseResp()
			return
		}
		var tagIDs []int64
		for _, tag := range tagMap {
			tagIDs = append(tagIDs, tag.ID)
		}
		postIDs, err := mysql.SelectArticleIDsByTagIDs(db, tagIDs)
		if err != nil {
			logger.Errorf("select post_id by tag_id error:[%v]", err)
			rsp.BaseResp = resp.NewFailBaseResp()
			return
		}
		if len(searchIDs) > 0 {
			// 取交集
			searchIDs = utils.IntersectInt64Slice(searchIDs, postIDs)
		} else {
			searchIDs = append(searchIDs, postIDs...)
		}
	}

	if searchByID && len(searchIDs) == 0 {
		return &api.GetPostListAPIResponse{
			Pagination: &api.Pagination{
				Page:    req.GetPage(),
				Limit:   req.GetLimit(),
				HasMore: false,
				Total:   &total,
			},
			BaseResp: resp.NewSuccessBaseResp(),
		}
	}

	// 搜索文章
	posts, err := mysql.SearchPostListByLimit(db, req.Keyword, searchIDs, req.Page, req.Limit)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return &api.GetPostListAPIResponse{
				Pagination: &api.Pagination{
					Page:    req.GetPage(),
					Limit:   req.GetLimit(),
					HasMore: false,
					Total:   &total,
				},
				BaseResp: resp.NewSuccessBaseResp(),
			}
		}
		logger.Errorf("search post by page error:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return
	}

	// 总数
	total, err = mysql.SelectSearchPostCount(db, req.Keyword, searchIDs)
	if err != nil {
		logger.Errorf("select post count error:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return
	}

	// 作者
	editorIDMap := map[int64]bool{}
	for _, post := range posts {
		editorIDMap[post.CreateUser] = true
	}
	var editorIDList []int64
	for id := range editorIDMap {
		editorIDList = append(editorIDList, id)
	}
	userMap, err := mysql.MSelectUserByIDs(db, editorIDList)
	if err != nil {
		logger.Errorf("select user error:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return
	}

	// 分类
	var postIDList []int64
	for _, post := range posts {
		postIDList = append(postIDList, post.ID)
	}
	categoryIDsMap, err := mysql.MSelectCategoryIDsByArticleIDs(db, postIDList)
	if err != nil {
		logger.Errorf("select categoryIDs error:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return
	}
	var categoryIDList []int64
	for _, ids := range categoryIDsMap {
		for _, id := range ids {
			categoryIDList = append(categoryIDList, id)
		}
	}
	categoryMap, err := mysql.MSelectCategoryByIDs(db, categoryIDList)
	if err != nil {
		logger.Errorf("select category error:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return
	}

	// 结果处理
	var postList []*api.PostListItem
	for _, post := range posts {
		var categories []string
		ids := categoryIDsMap[post.ID]
		for _, id := range ids {
			categories = append(categories, categoryMap[id].CategoryName)
		}
		userName := userMap[post.CreateUser].Nickname
		publishAt := int64(0)
		if post.PublishAt != nil {
			publishAt = post.PublishAt.Unix()
		}
		postList = append(postList, &api.PostListItem{
			ID:           post.ID,
			Title:        post.Title,
			CategoryList: categories,
			Editor:       userName,
			Status:       post.Status,
			PV:           post.PV,
			UpdateAt:     post.UpdateAt.Unix(),
			PublishAt:    publishAt,
		})
	}
	return &api.GetPostListAPIResponse{
		Pagination: &api.Pagination{
			Page:    req.GetPage(),
			Limit:   req.GetLimit(),
			HasMore: false,
			Total:   &total,
		},
		PostList: postList,
		BaseResp: resp.NewSuccessBaseResp(),
	}
}
