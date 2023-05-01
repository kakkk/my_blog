package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/storage"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreatePostAPI(ctx context.Context, req *api.CreatePostAPIRequest) (rsp *api.CreatePostAPIResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(
		logrus.Fields{
			"title": req.GetTitle(),
		},
	)
	rsp = &api.CreatePostAPIResponse{}
	var post *entity.Article

	// 开启事务
	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 校验分类是否存在
		err := checkCategoryIsExist(tx, req.CategoryList)
		if err != nil {
			if errors.Is(err, consts.ErrRecordNotFound) {
				rsp.BaseResp = resp.NewBaseResponse(common.RespCode_NotFound, "category not found")
				return fmt.Errorf("check category not exist: [%v]", err)
			}
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("check category exist error:[%v]", err)
		}

		// 创建文章
		userID, _ := utils.GetUserIDByCtx(ctx)
		var publishAt *time.Time
		if req.GetStatus() == common.ArticleStatus_PUBLISH {
			now := time.Now()
			publishAt = &now
		}
		post, err = mysql.CreateArticle(tx, &entity.Article{
			Title:       req.GetTitle(),
			Content:     req.GetContent(),
			ArticleType: common.ArticleType_Post,
			Status:      req.GetStatus(),
			CreateUser:  userID,
			PublishAt:   publishAt,
		})
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("create article error:[%v]", err)
		}

		// 处理标签
		err = processTags(tx, req.GetTags(), post.ID)
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("process tag error:[%v]", err)
		}

		// 处理分类
		err = processCategories(tx, req.GetCategoryList(), post.ID)
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("process categories error:[%v]", err)
		}

		// 提交事务
		return nil
	})

	if err != nil {
		logger.Errorf("create post fail:[%v]", err)
		return rsp
	}

	logger.Infof("create post success")
	// 重建缓存
	storage.GetPostOrderListStorage().Rebuild(ctx)
	return &api.CreatePostAPIResponse{
		ID:       post.ID,
		BaseResp: resp.NewSuccessBaseResp(),
	}
}

func GetPostAPI(ctx context.Context, req *api.GetPostAPIRequest) *api.GetPostAPIResponse {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req,
	})
	db := mysql.GetDB(ctx)
	failResp := &api.GetPostAPIResponse{
		BaseResp: resp.NewBaseResponse(common.RespCode_NotFound, "not found"),
	}

	// 获取post
	post, err := mysql.SelectArticleByID(db, req.GetID())
	if err != nil {
		if err == consts.ErrRecordNotFound {
			logger.Warnf("post not found, post_id:[%v]", req.GetID())
			return failResp
		}
		logger.Errorf("select post error:[%v]", err)
		return failResp
	}

	// 获取分类
	categoryIDs, err := mysql.SelectCategoryIDsByArticleID(db, post.ID)
	if err != nil {
		logger.Errorf("select category_id error:[%v]", err)
		return failResp
	}
	categoryMap, err := mysql.MSelectCategoryByIDs(db, categoryIDs)
	if err != nil {
		logger.Errorf("select category error:[%v]", err)
		return failResp
	}
	var categoryList []*api.CategoriesItem
	for _, category := range categoryMap {
		categoryList = append(categoryList, &api.CategoriesItem{
			ID:   category.ID,
			Name: category.CategoryName,
		})
	}

	// 标签
	tagList, err := mysql.SelectTagListByArticleID(db, req.GetID())
	if err != nil {
		logger.Errorf("select tag list error:[%v]", err)
		return failResp
	}

	return &api.GetPostAPIResponse{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		Status:       post.Status,
		CategoryList: categoryList,
		Tags:         tagList,
		BaseResp:     resp.NewSuccessBaseResp(),
	}
}

func UpdatePostAPI(ctx context.Context, req *api.UpdatePostAPIRequest) (rsp *api.CommonResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req.GetID(),
		"title":   req.GetTitle(),
	})
	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 校验分类是否存在
		err := checkCategoryIsExist(tx, req.CategoryList)
		if err != nil {
			if errors.Is(err, consts.ErrRecordNotFound) {
				rsp.BaseResp = resp.NewBaseResponse(common.RespCode_NotFound, "category not found")
				return fmt.Errorf("check category not exist: [%v]", err)
			}
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("check category exist error:[%v]", err)
		}

		// 更新文章
		err = mysql.UpdateArticleByID(tx, req.GetID(), &entity.Article{
			ID:      0,
			Title:   req.GetTitle(),
			Content: req.GetContent(),
		})
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("update post error:[%v]", err)
		}

		// 处理标签
		err = processTags(tx, req.GetTags(), req.GetID())
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("process tag error:[%v]", err)
		}

		// 处理分类
		err = processCategories(tx, req.GetCategoryList(), req.GetID())
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("process categories error:[%v]", err)
		}

		return nil
	})

	if err != nil {
		logger.Errorf("update post fail:[%v]", err)
		return rsp
	}

	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
}

// 检查分类是否存在
func checkCategoryIsExist(tx *gorm.DB, categoryIDs []int64) error {
	categoryMap, err := mysql.MSelectCategoryByIDs(tx, categoryIDs)
	if err != nil {
		return fmt.Errorf("select category error:[%v]", err)
	}
	for _, id := range categoryIDs {
		_, ok := categoryMap[id]
		if !ok {
			return fmt.Errorf("check category is exist error:[%w]", consts.ErrRecordNotFound)
		}
	}
	return nil
}

// 标签处理
func processTags(tx *gorm.DB, tags []string, postID int64) error {
	// 先清除标签（更新操作会用到）
	err := mysql.DeleteArticleTagRelationByArticleID(tx, postID)
	if err != nil {
		return fmt.Errorf("delete tag error:[%v]", err)
	}

	// 获取存在的标签
	tagMap, err := mysql.MSelectTagByName(tx, tags)
	if err != nil {
		return fmt.Errorf("select tags error:[%v]", err)
	}
	notExistTag := make([]*entity.Tag, 0)    // 不存在
	needRestoreTagIDs := make([]int64, 0)    // 需要恢复
	tagIDList := make([]int64, 0, len(tags)) // 所有标签ID
	for _, name := range tags {
		tag, ok := tagMap[name]
		if !ok {
			notExistTag = append(notExistTag, &entity.Tag{TagName: name})
			continue
		}
		if tag.DeleteFlag == common.DeleteFlag_Delete {
			needRestoreTagIDs = append(needRestoreTagIDs, tag.ID)
		}
		tagIDList = append(tagIDList, tag.ID)
	}

	// 创建不存在的标签
	createdTags, err := mysql.BatchCreateTags(tx, notExistTag)
	if err != nil {
		return fmt.Errorf("create tags error:[%v]", err)
	}
	for _, tag := range createdTags {
		tagIDList = append(tagIDList, tag.ID)
	}

	// 恢复被删除的标签
	err = mysql.RestoreTagByIDs(tx, needRestoreTagIDs)
	if err != nil {
		return fmt.Errorf("restore tags error:[%v]", err)
	}

	// 文章添加标签
	var tagRelation []*entity.ArticleTag
	for _, id := range tagIDList {
		tagRelation = append(tagRelation, &entity.ArticleTag{
			PostID: postID,
			TagID:  id,
		})
	}
	err = mysql.UpsertArticleTagRelation(tx, tagRelation)
	if err != nil {
		return fmt.Errorf("upsert article_tag error:[%v]", err)
	}

	// 刷新标签状态
	err = mysql.RefreshTagUpdateTimeByIDs(tx, tagIDList)
	if err != nil {
		return fmt.Errorf("refresh tag update_time error:[%v]", err)
	}

	return nil
}

func processCategories(tx *gorm.DB, categoryIDs []int64, postID int64) error {
	// 解除所有关联
	err := mysql.DeleteArticleCategoryRelationByArticleID(tx, postID)
	if err != nil {
		return fmt.Errorf("delete article category relation error:[%v]", err)
	}

	// 添加分类
	var categoryRelation []*entity.ArticleCategory
	for _, id := range categoryIDs {
		categoryRelation = append(categoryRelation, &entity.ArticleCategory{
			PostID:     postID,
			CategoryID: id,
		})
	}
	// 默认分类
	if len(categoryIDs) == 0 {
		defaultID, err := mysql.SelectDefaultCategoryID(tx)
		if err != nil {
			return fmt.Errorf("select default category_id error:[%v]", err)
		}
		categoryRelation = append(categoryRelation, &entity.ArticleCategory{
			PostID:     postID,
			CategoryID: defaultID,
		})
	}
	err = mysql.UpsertArticleCategoryRelation(tx, categoryRelation)
	if err != nil {
		return fmt.Errorf("upsert article_category error:[%v]", err)
	}

	return nil
}

func UpdatePostStatusAPI(ctx context.Context, req *api.UpdatePostStatusAPIRequest) (rsp *api.CommonResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req.GetID(),
		"status":  req.GetStatus(),
	})
	rsp = &api.CommonResponse{}

	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		post, err := mysql.SelectArticleByID(tx, req.GetID())
		if err != nil {
			if errors.Is(err, consts.ErrRecordNotFound) {
				rsp.BaseResp = resp.NewBaseResponse(common.RespCode_NotFound, "post not found")
				return fmt.Errorf("post not found")
			}
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("select article by id error:[%v]", err)
		}

		// 新发布文章
		if req.GetStatus() == common.ArticleStatus_PUBLISH && post.PublishAt == nil {
			publishAt := time.Now()
			err = mysql.UpdateArticlePublishAtByID(tx, req.GetID(), &publishAt)
			if err != nil {
				rsp.BaseResp = resp.NewFailBaseResp()
				return fmt.Errorf("update article publish_at by id error:[%v]", err)
			}
		}

		// 更新状态
		err = mysql.UpdateArticleStatusByID(tx, req.GetID(), req.GetStatus())
		if err != nil {
			rsp.BaseResp = resp.NewFailBaseResp()
			return fmt.Errorf("update article status by id error:[%v]", err)
		}

		// 提交事务
		return nil
	})

	if err != nil {
		logger.Errorf("update article status fail:[%v]", err)
		return rsp
	}

	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}
}

func DeletePostAPI(ctx context.Context, req *api.DeletePostAPIRequest) (rsp *api.CommonResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"post_id": req.GetID(),
	})
	rsp = &api.CommonResponse{}
	err := mysql.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		err := mysql.DeleteArticleByID(tx, req.GetID())
		if err != nil {
			return fmt.Errorf("delete article error:[%v]", err)
		}
		err = mysql.DeleteArticleCategoryRelationByArticleID(tx, req.GetID())
		if err != nil {
			return fmt.Errorf("delete article_category relation error:[%v]", err)
		}
		err = mysql.DeleteArticleTagRelationByArticleID(tx, req.GetID())
		if err != nil {
			return fmt.Errorf("delete article_tag relation error:[%v]", err)
		}
		return nil
	})

	if err != nil {
		logger.Errorf("delete post fail:[%v]", err)
		rsp.BaseResp = resp.NewFailBaseResp()
		return
	}

	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}

}
