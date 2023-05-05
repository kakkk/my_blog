package service

import (
	"context"
	"my_blog/biz/common/consts"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/storage"

	"github.com/sirupsen/logrus"
)

func CreateCategoryAPI(ctx context.Context, req *api.CreateCategoryAPIRequest) *api.CommonResponse {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"category_name": req.GetName(),
		"slug":          req.GetSlug(),
	})
	tx := mysql.GetDB(ctx).Begin()

	// 创建分类
	category, err := mysql.CreateCategory(tx, &entity.Category{
		CategoryName: req.GetName(),
		Slug:         req.GetSlug(),
	})
	if err != nil {
		logger.Errorf("create category error:[%v]", err)
		tx.Rollback()
		if err == consts.ErrHasExist {
			return &api.CommonResponse{
				BaseResp: resp.NewBaseResponse(common.RespCode_HasExist, "has exist"),
			}
		}
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	// 更新排序
	order, err := mysql.SelectCategoryOrder(tx)
	if err != nil {
		logger.Errorf("select category order error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}
	order = append(order, category.ID)
	err = mysql.UpdateCategoryOrder(tx, order)
	if err != nil {
		logger.Errorf("update category order error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("commit transaction error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	logger.Infof("create category success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
}

func UpdateCategoryAPI(ctx context.Context, req *api.UpdateCategoryAPIRequest) *api.CommonResponse {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"id":            req.GetID(),
		"category_name": req.GetName(),
		"slug":          req.GetSlug(),
	})
	err := mysql.UpdateCategoryByID(mysql.GetDB(ctx), req.GetID(), &entity.Category{
		CategoryName: req.GetName(),
		Slug:         req.GetSlug(),
	})
	if err != nil {
		logger.Errorf("update category error:[%v]", err)
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}
	logger.Infof("update category success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
}

func DeleteCategoryAPI(ctx context.Context, req *api.DeleteCategoryAPIRequest) *api.CommonResponse {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"id": req.GetID(),
	})
	tx := mysql.GetDB(ctx).Begin()

	// 删除分类
	err := mysql.DeleteCategoryByID(tx, req.GetID())
	if err != nil {
		tx.Rollback()
		logger.Errorf("delete category error:[%v]", err)
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	// 更新排序
	beforeOrder, err := mysql.SelectCategoryOrder(tx)
	if err != nil {
		logger.Errorf("select category order error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}
	afterOrder := make([]int64, 0, len(beforeOrder)-1)
	for _, id := range beforeOrder {
		if id == req.GetID() {
			continue
		}
		afterOrder = append(afterOrder, id)
	}
	err = mysql.UpdateCategoryOrder(tx, afterOrder)
	if err != nil {
		logger.Errorf("update category order error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("commit transaction error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	logger.Infof("delete category success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
}

func GetCategoryListAPI(ctx context.Context) *api.GetCategoryListAPIResponse {
	logger := log.GetLoggerWithCtx(ctx)
	list, err := storage.GetCategoryListStorage().GetFromDB(ctx)
	if err != nil {
		logger.Errorf("get category list error:[%v]", err)
		return &api.GetCategoryListAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}
	}

	logger.Infof("get category list success")
	return &api.GetCategoryListAPIResponse{
		CategoryList: list.ToAPICategoryListModel(),
		BaseResp:     resp.NewSuccessBaseResp(),
	}
}

// TODO
func UpdateCategoryOrderAPI(ctx context.Context, req *api.UpdateCategoryOrderAPIRequest) *api.CommonResponse {
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
}
