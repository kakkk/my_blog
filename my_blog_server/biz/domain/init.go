package domain

import (
	"context"

	"my_blog/biz/domain/repo"
	"my_blog/biz/domain/repo/cache"
	"my_blog/biz/domain/service"
)

func MustInit() {
	ctx := context.Background()
	repo.InitRepo(ctx)
	MustInitIndex(ctx)
}

func MustInitIndex(ctx context.Context) {
	cache.InitArchivesStorage()
	err := service.GetIndexService().InitIndex(ctx)
	if err != nil {
		panic(err)
	}
}
