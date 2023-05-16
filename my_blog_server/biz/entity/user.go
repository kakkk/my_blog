package entity

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
