package storage

import (
	"fmt"

	"my_blog/biz/common/consts"
	"my_blog/biz/components/cachex"
)

func parseSqlError[T any](err error) (T, error) {
	var zero T
	if err == consts.ErrRecordNotFound {
		return zero, cachex.ErrNotFound
	}
	return zero, fmt.Errorf("sql error: %w", err)
}
