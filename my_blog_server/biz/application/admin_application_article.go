package application

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo"
	"my_blog/biz/domain/service"
	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/infra/session"
)

func (a *AdminApplication) CreatePost(ctx context.Context, req *api.CreatePostAPIRequest) (*api.CreatePostAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(
		logrus.Fields{
			"title": req.GetTitle(),
		},
	)
	userID, _ := session.GetUserIDByCtx(ctx)
	var publishAt *time.Time
	if req.GetStatus() == common.ArticleStatus_PUBLISH {
		now := time.Now()
		publishAt = &now
	}
	article := &dto.Article{
		Title:        req.GetTitle(),
		Content:      req.GetContent(),
		ArticleType:  common.ArticleType_Post,
		Status:       req.GetStatus(),
		CreateUserID: userID,
		PublishAt:    publishAt,
	}
	categories := req.GetCategoryList()
	if len(categories) == 0 {
		categories = append(categories, repo.GetContentRepo().Cache().GetDefaultCategoryID())
	}
	// 创建文章
	post, err := service.GetContentService().CreateArticle(ctx, article, req.GetCategoryList(), req.GetTags())
	if err != nil {
		logger.Errorf("create article fail: %v", err)
		return nil, fmt.Errorf("create article fail: %w", err)
	}
	// 更新搜索索引
	err = service.GetIndexService().IndexArticle(&entity.ArticleData{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
	})
	if err != nil {
		logger.Errorf("index article fail: %v", err)
		return nil, fmt.Errorf("index article fail: %w", err)
	}
	// 更新归档
	service.GetIndexService().RefreshArchives(ctx, nil)
	// 返回
	return &api.CreatePostAPIResponse{
		ID:       post.ID,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) GetPostByID(ctx context.Context, req *api.GetPostAPIRequest) (*api.GetPostAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(
		logrus.Fields{
			"id": req.GetID(),
		},
	)
	article, err := service.GetContentService().GetArticleFromDB(ctx, req.GetID())
	if err != nil {
		logger.Errorf("get article from db fail: %v", err)
		return nil, fmt.Errorf("get article from db fail: %v", err)
	}

	return &api.GetPostAPIResponse{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		Status:       article.Status,
		CategoryList: dto.Categories(article.Categories).ToCategoriesItems(),
		Tags:         dto.Tags(article.Tags).ToStringList(),
		BaseResp:     resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) UpdatePost(ctx context.Context, req *api.UpdatePostAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req.GetID(),
		"title":   req.GetTitle(),
	})
	article := &dto.Article{
		ID:      req.GetID(),
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}
	err := service.GetContentService().UpdateArticle(ctx, article, req.GetCategoryList(), req.GetTags())
	if err != nil {
		logger.Errorf("create article fail: %v", err)
		return nil, fmt.Errorf("create article fail: %w", err)
	}
	logger.Infof("update post success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

// UpdatePostStatus 更新文章状态
func (a *AdminApplication) UpdatePostStatus(ctx context.Context, req *api.UpdatePostStatusAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req.GetID(),
		"status":  req.GetStatus(),
	})
	err := service.GetContentService().UpdateArticleStatus(ctx, req.GetID(), req.GetStatus())
	if err != nil {
		return nil, fmt.Errorf("update article status fail:[%w]", err)
	}
	logger.Infof("update post status success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) DeletePost(ctx context.Context, req *api.DeletePostAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req.GetID(),
	})
	article := entity.NewArticleByDTO(&dto.Article{ID: req.GetID()}, nil, nil)
	err := article.Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("delete post fail: %w", err)
	}
	logger.Infof("delete post success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) GetPostList(ctx context.Context, req *api.GetPostListAPIRequest) (*api.GetPostListAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	articles, page, err := service.GetContentService().SearchArticleList(ctx, req.Keyword, req.GetCategories(), req.GetTags(), req.Page, req.Limit)
	if err != nil {
		logger.Errorf("search fail: %v", err)
		return nil, fmt.Errorf("search fail: %w", err)
	}
	logger.Infof("search success")
	return &api.GetPostListAPIResponse{
		Pagination: page.ToRespPagination(req.GetPage(), req.GetLimit()),
		PostList:   dto.Articles(articles).ToRespList(),
		BaseResp:   resp.NewSuccessBaseResp(),
	}, nil
}
