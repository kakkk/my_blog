package service

import (
	"context"

	"my_blog/biz/common/consts"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/model/blog/api"
	"my_blog/biz/model/blog/common"
	"my_blog/biz/repository/mysql"
)

func LoginAPI(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)
	user, err := mysql.SelectUserByUsername(ctx, mysql.GetDB(), req.GetUsername())
	if err != nil {
		if err == consts.ErrRecordNotFound {
			logger.Warnf("login fail: user not found, username:[%v], password:[%v]", req.GetUsername(), req.GetPassword())
			return &api.LoginResponse{
				BaseResp: resp.NewBaseResponse(common.RespCode_LoginFail, "login fail"),
			}, nil
		}
	}
	pass := utils.CompareHashAndPassword(user.PwdHash, req.GetPassword())
	if !pass {
		logger.Warnf("login fail: password incorrect, username:[%v], password:[%v]", req.GetUsername(), req.GetPassword())
		return &api.LoginResponse{
			BaseResp: resp.NewBaseResponse(common.RespCode_LoginFail, "login fail"),
		}, nil
	}

	return &api.LoginResponse{
		UserID:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		BaseResp: resp.NewBaseResponse(common.RespCode_Success, ""),
	}, nil
}

func GetUserInfoAPI(ctx context.Context) (*api.GetUserInfoAPIResponse, error) {
	logger := log.GetLoggerWithCtx(ctx)

	userID, err := utils.GetUserIDByCtx(ctx)
	if err != nil {
		return &api.GetUserInfoAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	user, err := mysql.SelectUserByID(ctx, mysql.GetDB(), userID)
	if err != nil {
		logger.Errorf("get user by id error: %v", err)
		return &api.GetUserInfoAPIResponse{
			BaseResp: resp.NewFailBaseResp(),
		}, nil
	}

	return &api.GetUserInfoAPIResponse{
		UserID:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		BaseResp: resp.NewBaseResponse(common.RespCode_Success, ""),
	}, nil
}
