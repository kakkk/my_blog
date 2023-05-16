package storage

import (
	"context"
	"fmt"

	"my_blog/biz/common/log"
	"my_blog/biz/components/cachex"
)

func InitStorage() error {
	ctx := context.Background()
	cachex.SetLogger(log.NewCacheXLogger())
	err := initArticleEntityStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostOrderListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostPrevNextStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostMetaStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initUserEntityStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategoryEntityStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostTagListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initPostCategoryListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategoryListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategoryPostListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initCategorySlugIDStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initTagListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initTagNameIDStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	err = initTagPostListStorage(ctx)
	if err != nil {
		return fmt.Errorf("init storage error: [%w]", err)
	}
	initArchivesStorage()
	return nil
}
