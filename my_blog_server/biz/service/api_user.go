package service

import (
	"context"

	"my_blog/biz/common/resp"
	"my_blog/biz/consts"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/log"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/infra/session"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
	"my_blog/biz/repository/mysql"
)

func LoginAPI(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	user, err := mysql.SelectUserByUsername(mysql2.GetDB(ctx), req.GetUsername())
	if err != nil {
		if err == consts.ErrRecordNotFound {
			logger.Warnf("login fail: user not found, username:[%v], password:[%v]", req.GetUsername(), req.GetPassword())
			return &api.LoginResponse{
				BaseResp: resp.NewBaseResponse(common.RespCode_LoginFail, "login fail"),
			}, nil
		}
	}
	pass := misc.CompareHashAndPassword(user.PwdHash, req.GetPassword())
	if !pass {
		logger.Warnf("login fail: password incorrect, username:[%v], password:[%v]", req.GetUsername(), req.GetPassword())
		return &api.LoginResponse{
			BaseResp: resp.NewBaseResponse(common.RespCode_LoginFail, "login fail"),
		}, nil
	}

	err = session.SetUserID(ctx, user.ID)
	if err != nil {
		logger.Errorf("set session error:[%v]", err)
		return &api.LoginResponse{
			BaseResp: resp.NewBaseResponse(common.RespCode_LoginFail, "login fail"),
		}, nil
	}

	return &api.LoginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}

func GetUserInfoAPI(ctx context.Context) (*api.GetUserInfoAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)

	userID, err := session.GetUserIDByCtx(ctx)
	if err != nil {
		return &api.GetUserInfoAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	user, err := mysql.SelectUserByID(mysql2.GetDB(ctx), userID)
	if err != nil {
		logger.Errorf("get user by id error: %v", err)
		return &api.GetUserInfoAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	return &api.GetUserInfoAPIResponse{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		BaseResp: resp.NewSuccessBaseResp(),
	}, nil
}
