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
)

func CreateTagAPI(ctx context.Context, req *api.CreateTagAPIRequest) (*api.CreateTagAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)

	tag, err := mysql.CreateTag(mysql.GetDB(ctx), &entity.Tag{
		TagName: req.Name,
	})
	if err != nil {
		if err == consts.ErrHasExist {
			logger.Warnf("tag has exist, tag_name:[%v]", req.GetName())
			return &api.CreateTagAPIResponse{
				BaseResp: resp.NewBaseResponse(common.RespCode_HasExist, "has exist"),
			}, nil
		}
		logger.Errorf("create tag fail, error:[%v]", err)
		return &api.CreateTagAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	logger.Infof("create tag success, tag_name:[%v]", req.GetName())
	return &api.CreateTagAPIResponse{
		ID:       tag.ID,
		Name:     tag.TagName,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func UpdateTagAPI(ctx context.Context, req *api.UpdateTagAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	err := mysql.UpdateTagByID(mysql.GetDB(ctx), req.GetID(), &entity.Tag{
		TagName: req.Name,
	})
	if err != nil {
		if err == consts.ErrHasExist {
			logger.Warnf("tag has exist, tag_name:[%v]", req.GetName())
			return &api.CommonResponse{
				BaseResp: resp.NewBaseResponse(common.RespCode_HasExist, "has exist"),
			}, nil
		}
		logger.Errorf("update tag fail, error:[%v]", err)
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	logger.Infof("update tag success, tag_name:[%v]", req.GetName())
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func DeleteTagAPI(ctx context.Context, req *api.DeleteTagAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithField("tag_id", req.GetID())

	tx := mysql.GetDB(ctx).Begin()

	err := mysql.DeleteTagByID(tx, req.GetID())
	if err != nil {
		logger.Errorf("delete tag fail, error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	err = mysql.DeleteArticleTagRelationByTagID(tx, req.GetID())
	if err != nil {
		logger.Errorf("delete article_tag fail, error:[%v]", err)
		tx.Rollback()
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	err = tx.Commit().Error
	if err != nil {
		logger.Errorf("commit transaction fail, error:[%v]", err)
		return &api.CommonResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	logger.Infof("delete tag success")
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func GetTagListAPI(ctx context.Context, req *api.GetTagListAPIRequest) (*api.GetTagListAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	result, err := mysql.GetTagListByPage(mysql.GetDB(ctx), req.Keyword, req.Page, req.Limit)
	if err != nil {
		logger.Errorf("get tag list fail, error:[%v]", err)
		return &api.GetTagListAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}
	hasMore := len(result) < int(req.GetLimit())

	tagIDs := make([]int64, 0, len(result))
	for _, tag := range result {
		tagIDs = append(tagIDs, tag.ID)
	}
	counts, err := mysql.MGetTagArticleCountByTagIDs(mysql.GetDB(ctx), tagIDs)
	if err != nil {
		logger.Warnf("mget tag article count error:[%v]", err)
		return &api.GetTagListAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}
	tagList := make([]*api.TagListItem, 0, len(result))
	for _, tag := range result {
		tagList = append(tagList, &api.TagListItem{
			ID:    tag.ID,
			Name:  tag.TagName,
			Count: counts[tag.ID],
		})
	}
	allCount, err := mysql.GetAllTagCount(mysql.GetDB(ctx))
	if err != nil {
		logger.Warnf("get all tag count error:[%v]", err)
	}

	logger.Infof("get tag list success")
	return &api.GetTagListAPIResponse{
		Pagination: &api.Pagination{
			Page:    req.GetPage(),
			Limit:   req.GetLimit(),
			HasMore: hasMore,
			Total:   &allCount,
		},
		TagList:  tagList,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}
