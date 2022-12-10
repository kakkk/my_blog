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

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/sessions"
)

func LoginAPI(ctx context.Context, c *app.RequestContext, req api.LoginRequest) (*api.LoginResponse, error) {
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
	session := sessions.Default(c)
	session.Set("user_id", user.Id)
	err = session.Save()
	if err != nil {
		logger.Errorf("set session error:[%v]", err)
		return &api.LoginResponse{
			BaseResp: resp.NewBaseResponse(common.RespCode_Fail, "fail"),
		}, nil
	}
	return &api.LoginResponse{
		UserID:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		BaseResp: resp.NewBaseResponse(common.RespCode_Success, ""),
	}, nil
}
