package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/dto"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	cachex2 "my_blog/biz/infra/repository/cachex"
	mysql2 "my_blog/biz/infra/repository/mysql"
	"my_blog/biz/repository/mysql"
)

var commentStorage *CommentStorage

type CommentStorage struct {
	cacheX *cachex.CacheX[int64, *dto.Comment]
	expire time.Duration
}

func GetCommentStorageStorage() *CommentStorage {
	return commentStorage
}

func initCommentStorageStorage(ctx context.Context) error {
	cfg := config.GetStorageSettingByName("comment")
	cx, err := cachex2.NewCacheXBuilderByConfig[int64, *dto.Comment](ctx, cfg).
		SetGetRealData(commentStorageGetRealData).
		SetMGetRealData(commentStorageMGetRealData).
		Build()
	if err != nil {
		return fmt.Errorf("init cachex error: %w", err)
	}
	commentStorage = &CommentStorage{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
	return nil
}

func commentStorageGetRealData(ctx context.Context, id int64) (*dto.Comment, error) {
	commentEntity, err := mysql.SelectApprovalCommentByID(mysql2.GetDB(ctx), id)
	if err != nil {
		return parseSqlError[*dto.Comment](nil, err)
	}
	return dto.NewCommentWithEntity(commentEntity), nil
}

func commentStorageMGetRealData(ctx context.Context, ids []int64) (map[int64]*dto.Comment, error) {
	commentEntities, err := mysql.SelectApprovalCommentsByIDs(mysql2.GetDB(ctx), misc.SliceDeduplicate(ids))
	if err != nil {
		return parseSqlError[map[int64]*dto.Comment](nil, err)
	}
	comments := make(map[int64]*dto.Comment, len(commentEntities))
	for _, entity := range commentEntities {
		comments[entity.ID] = dto.NewCommentWithEntity(entity)
	}
	return comments, nil
}

func (c *CommentStorage) MGet(ctx context.Context, ids []int64) map[int64]*dto.Comment {
	return c.cacheX.MGet(ctx, ids, c.expire)
}

func (c *CommentStorage) Get(ctx context.Context, id int64) (*dto.Comment, error) {
	comment, ok := c.cacheX.Get(ctx, id, c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return comment, nil
}
