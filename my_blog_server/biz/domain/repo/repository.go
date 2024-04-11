package repo

import (
	"context"

	"my_blog/biz/domain/repo/impl"
	"my_blog/biz/domain/repo/interfaces"
)

var contentRepo = impl.ContentRepoImpl{}
var manageRepo = impl.ManageRepoImpl{}
var commentRepo = impl.CommentRepoImpl{}

func GetContentRepo() interfaces.ContentRepo {
	return contentRepo
}

func GetManageRepo() interfaces.ManageRepo {
	return manageRepo
}

func GetCommentRepo() interfaces.CommentRepo {
	return commentRepo
}

func InitRepo(ctx context.Context) {
	impl.InitContentCache(ctx)
	impl.MustInitCachex()
}
