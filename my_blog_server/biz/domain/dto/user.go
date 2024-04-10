package dto

import "my_blog/biz/infra/repository/model"

type User struct {
	ID       int64
	Username string
	Nickname string
	Avatar   string
	PwdHash  string
}

func DefaultUser() *User {
	// TODO 可配置
	return &User{
		ID:       0,
		Username: "kakkk",
		Nickname: "kakkk",
	}
}

func NewUserByModel(user *model.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		PwdHash:  user.PwdHash,
	}
}

func NewUserMapByModelMap(users map[int64]*model.User) map[int64]*User {
	result := make(map[int64]*User, len(users))
	for id, user := range users {
		result[id] = NewUserByModel(user)
	}
	return result
}
