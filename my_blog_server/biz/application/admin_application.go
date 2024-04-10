package application

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"my_blog/biz/domain/repo"
	"my_blog/biz/domain/service"
	"my_blog/biz/infra/pkg/log"
	"my_blog/biz/infra/pkg/resp"
	"my_blog/biz/infra/repository/mysql"
	"my_blog/biz/infra/session"

	"my_blog/biz/model/blog/api"
)

type AdminApplication struct{}

func GetAdminApplication() *AdminApplication {
	return &AdminApplication{}
}

func (*AdminApplication) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	logger := log.GetLoggerWithCtx(ctx).WithFields(logrus.Fields{
		"id": req.GetUsername(),
	})

	user, err := service.GetManageService().Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		logger.Errorf("delete category fail: %v", err)
		return nil, err
	}

	return &api.LoginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func (*AdminApplication) GetUserInfo(ctx context.Context) (*api.GetUserInfoAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)

	userID, err := session.GetUserIDByCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user id by ctx fail: %v", err)
	}

	user, err := repo.GetManageRepo().SelectUserByID(mysql.GetDB(ctx), userID)
	if err != nil {
		logger.Errorf("get user by id error: %v", err)
		return nil, fmt.Errorf("get user by id error: %w", err)
	}

	return &api.GetUserInfoAPIResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}
