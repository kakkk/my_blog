package storage

import (
	"time"

	"my_blog/biz/model/blog/page"
)

const ArchivesStorageTTL = 30 // 过期时间，单位：分钟

var archivesStorage *ArchivesStorage

// 直接使用本地缓存即可
type ArchivesStorage struct {
	cache    []*page.ArchiveByYear
	updateAt int64
}

func initArchivesStorage() {
	archivesStorage = &ArchivesStorage{
		cache:    []*page.ArchiveByYear{},
		updateAt: 0, // 初始化的UpdateAt为0，一定会返回过期
	}
}

func GetArchivesStorage() *ArchivesStorage {
	return archivesStorage
}

func (s *ArchivesStorage) Get() ([]*page.ArchiveByYear, bool) {
	// 过期，但返回旧值，用作兜底
	if s.updateAt+ArchivesStorageTTL*60 < time.Now().Unix() {
		return s.cache, true
	}
	return s.cache, false
}

func (s *ArchivesStorage) Set(archives []*page.ArchiveByYear) {
	s.updateAt = time.Now().Unix()
	s.cache = archives
}
