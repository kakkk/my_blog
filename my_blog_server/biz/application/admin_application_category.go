package application

import (
	"context"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/api"
)

func (a *AdminApplication) CreateCategory(ctx context.Context, req *api.CreateCategoryAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"category_name": req.GetName(),
		"slug":          req.GetSlug(),
	})
	category := &entity.Category{
		CategoryName: req.GetName(),
		Slug:         req.GetSlug(),
	}
	err := category.Create(ctx)
	if err != nil {
		logger.Errorf("save category fail: %v", err)
		return nil, err
	}
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) UpdateCategory(ctx context.Context, req *api.UpdateCategoryAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"id":            req.GetID(),
		"category_name": req.GetName(),
		"slug":          req.GetSlug(),
	})
	category := &entity.Category{
		ID:           req.GetID(),
		CategoryName: req.GetName(),
		Slug:         req.GetSlug(),
	}
	err := category.Update(ctx)
	if err != nil {
		logger.Errorf("update category fail: %v", err)
		return nil, err
	}
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) DeleteCategory(ctx context.Context, req *api.DeleteCategoryAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"id": req.GetID(),
	})
	category := &entity.Category{
		ID: req.GetID(),
	}
	err := category.Delete(ctx)
	if err != nil {
		logger.Errorf("delete category fail: %v", err)
		return nil, err
	}
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) GetCategoryList(ctx context.Context) (*api.GetCategoryListAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	list, err := repo.GetContentRepo().GetCategoryList(mysql.GetDB(ctx), false)
	if err != nil {
		logger.Errorf("get category list error:[%v]", err)
		return nil, err
	}

	logger.Infof("get category list success")
	return &api.GetCategoryListAPIResponse{
		CategoryList: dto.Categories(list).ToCategoryListItems(),
		BaseResp:     resp.NewSuccessBaseResp(),
	}, nil
}
