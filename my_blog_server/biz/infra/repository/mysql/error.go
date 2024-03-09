package mysql

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"my_blog/biz/consts"
)

func ParseError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return consts.ErrRecordNotFound
	}
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return fmt.Errorf("db err: [%v]", err)
	}
	switch mysqlErr.Number {
	case 1062:
		return consts.ErrHasExist
	}
	return fmt.Errorf("db err: [%v]", err)
}
