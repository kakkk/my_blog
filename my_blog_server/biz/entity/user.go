package entity

import "encoding/json"

type User struct {
	ID       int64  `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Nickname string `gorm:"column:nickname"`
	Avatar   string `gorm:"column:avatar"`
	PwdHash  string `gorm:"column:pwd_hash"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) Serialize() string {
	bytes, _ := json.Marshal(u)
	return string(bytes)
}

func (u *User) Deserialize(str string) (*User, error) {
	user := &User{}
	err := json.Unmarshal([]byte(str), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
