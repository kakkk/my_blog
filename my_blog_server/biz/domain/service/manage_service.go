package service

import (
	"context"
	"errors"
	"fmt"

	"my_blog/biz/consts"
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo"
	"my_blog/biz/infra/repository/mysql"
)

type ManageService struct{}

func GetManageService() ManageService {
	return ManageService{}
}

func (svc ManageService) Login(ctx context.Context, username string, pwd string) (*dto.User, error) {
	userDTO, err := repo.GetManageRepo().SelectUserByUsername(mysql.GetDB(ctx), username)
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, consts.ErrLoginFail
		}
		return nil, err
	}
	user := entity.NewUserByDTO(userDTO)
	ok := user.ComparePwhHash(pwd)
	if !ok {
		return nil, consts.ErrLoginFail
	}
	err = user.SessionSave(ctx)
	if !ok {
		return nil, fmt.Errorf("save session fail: %w", err)
	}
	return userDTO, nil
}
