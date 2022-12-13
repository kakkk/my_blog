package mysql

import (
	"context"
	"fmt"
	"time"

	"my_blog/biz/common/config"
	"my_blog/biz/common/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// InitMySQL 初始化MySQL
func InitMySQL() error {
	cfg := config.GetMySQLConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.NewGORMLogger(),
	})
	if err != nil {
		return err
	}
	db = gormDB
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

func GetDB(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}
