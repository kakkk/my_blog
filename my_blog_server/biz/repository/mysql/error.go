package mysql

import (
	"errors"
	"fmt"

	"my_blog/biz/common/consts"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func parseError(err error) error {
	if err == gorm.ErrRecordNotFound {
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
