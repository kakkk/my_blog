package storage

import (
	"errors"
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

func parseCacheXError[T any](val T, err error) (T, error) {
	var zero T
	if errors.Is(err, cachex.ErrNotFound) {
		return zero, consts.ErrRecordNotFound
	}
	return zero, fmt.Errorf("cachex error:[%v]", err)
}
