package impl

import "my_blog/biz/domain/repo/interfaces"

type ContentRepoImpl struct{}

func (c ContentRepoImpl) Cache() interfaces.ContentCache {
	return &ContentCacheImpl{}
}
