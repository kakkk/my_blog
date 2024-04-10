package cache

import (
	"sync/atomic"
	"time"

	"my_blog/biz/hertz_gen/blog/page"
)

const ArchivesCacheTTL = 1440 // 过期时间，单位：分钟

var archives *Archives

// 直接使用本地缓存即可
type Archives struct {
	cache    atomic.Value
	updateAt int64
}

func InitArchivesStorage() {
	archives = &Archives{
		cache:    atomic.Value{},
		updateAt: 0,
	}
	archives.cache.Store([]*page.ArchiveByYear{})
}

func GetArchivesStorage() *Archives {
	return archives
}

func (s *Archives) Get() ([]*page.ArchiveByYear, bool) {
	// 过期，但返回旧值，用作兜底
	if s.updateAt+ArchivesCacheTTL*60 < time.Now().Unix() {
		return s.cache.Load().([]*page.ArchiveByYear), true
	}
	return s.cache.Load().([]*page.ArchiveByYear), false
}

func (s *Archives) Set(a []*page.ArchiveByYear) {
	s.updateAt = time.Now().Unix()
	archives.cache.Store(a)
}
