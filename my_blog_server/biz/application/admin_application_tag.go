package application

import (
	"context"

	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/model/blog/api"
)

func (a *AdminApplication) CreateTag(ctx context.Context, req *api.CreateTagAPIRequest) (*api.CreateTagAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithField("tag_name", req.GetName())
	tag := entity.NewTagByDTO(&dto.Tag{
		TagName: req.GetName(),
	})
	err := tag.Create(ctx)
	if err != nil {
		return nil, err
	}
	logger.Infof("create tag success, id:%v", tag.ID)
	return &api.CreateTagAPIResponse{
		ID:       tag.ID,
		Name:     tag.TagName,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) UpdateTag(ctx context.Context, req *api.UpdateTagAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithField("tag_id", req.GetID())
	tag := entity.NewTagByDTO(&dto.Tag{
		TagName: req.GetName(),
	})
	err := tag.Update(ctx)
	if err != nil {
		return nil, err
	}
	logger.Infof("update tag success, id:%v", tag.ID)
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) DeleteTag(ctx context.Context, req *api.DeleteTagAPIRequest) (*api.CommonResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithField("tag_id", req.GetID())
	tag := entity.NewTagByID(req.GetID())
	err := tag.Delete(ctx)
	if err != nil {
		return nil, err
	}
	logger.Infof("delete tag success, id:%v", tag.ID)
	return &api.CommonResponse{
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (a *AdminApplication) GetTagList(ctx context.Context, req *api.GetTagListAPIRequest) (*api.GetTagListAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	list, count, err := repo.GetContentRepo().GetTagListByPage(mysql.GetDB(ctx), req.Keyword, req.Page, req.Limit)
	if err != nil {
		logger.Errorf("get category list error:[%v]", err)
		return nil, err
	}
	hasMore := len(list) < int(req.GetLimit())
	logger.Infof("get category list success")
	return &api.GetTagListAPIResponse{
		Pagination: &api.Pagination{
			Page:    req.GetPage(),
			Limit:   req.GetLimit(),
			HasMore: hasMore,
			Total:   &count,
		},
		TagList: dto.Tags(list).ToTagList(),

		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}
