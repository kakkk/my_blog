package entity

import (
	"context"

	"my_blog/biz/domain/dto"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/session"
)

type User struct {
	ID       int64
	Username string
	Nickname string
	Avatar   string
	PwdHash  string
}

// ComparePwhHash 比较密码
func (u *User) ComparePwhHash(pwd string) bool {
	return misc.CompareHashAndPassword(u.PwdHash, pwd)
}

// SessionSave 保存session
func (u *User) SessionSave(ctx context.Context) error {
	return session.SetUserID(ctx, u.ID)
}

func NewUserByDTO(user *dto.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		PwdHash:  user.PwdHash,
	}
}
