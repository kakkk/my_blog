package cache

import "sync/atomic"

var defaultCategoryID atomic.Int64

func GetDefaultCategoryID() int64 {
	return defaultCategoryID.Load()
}

func SetDefaultCategoryID(id int64) {
	defaultCategoryID.Store(id)
}
