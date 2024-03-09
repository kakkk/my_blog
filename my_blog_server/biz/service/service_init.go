package service

import (
	"context"
	"fmt"

	"my_blog/biz/infra/pkg/log"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/index"
	"my_blog/biz/repository/mysql"
)

// service初始化
func InitService() error {
	ctx := context.Background()
	logger := log.GetLoggerWithCtx(ctx)

	// 加载全量文章数据
	postFromDB, err := mysql.SelectAllPublishedPostWithBatch(mysql2.GetDB(ctx))
	if err != nil {
		return fmt.Errorf("load all post from db error:[%v]", err)
	}
	logger.Info("load all post from db success")

	// 加载文章归档
	RefreshArchives(ctx, postFromDB)
	logger.Info("load archives cache success")

	// 索引全量文章
	err = index.GetArticleIndex().MIndexArticle(postFromDB)
	if err != nil {
		return fmt.Errorf("index all post error:[%v]", err)
	}
	logger.Info("index all post success")

	return nil
}
