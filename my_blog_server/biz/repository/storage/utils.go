package storage

import (
	"errors"
	"fmt"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
)

func parseSqlError[T any](val T, err error) (T, error) {
	var zero T
	if errors.Is(err, consts.ErrRecordNotFound) {
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
