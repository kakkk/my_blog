package cachex

import "errors"

var (
	ErrNotFound         = errors.New("data not found")
	ErrCacheError       = errors.New("cache error")
	ErrGetRealDataError = errors.New("get real data fail")
)
