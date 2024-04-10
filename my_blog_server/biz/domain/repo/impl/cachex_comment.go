package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/kakkk/cachex"

	"my_blog/biz/consts"
	"my_blog/biz/domain/dto"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	infraCachex "my_blog/biz/infra/repository/cachex"
	"my_blog/biz/infra/repository/mysql"
)

type CommentCachex struct {
	cacheX *cachex.CacheX[int64, *dto.Comment]
	expire time.Duration
}

var commentCachex *CommentCachex

func initCommentCachex() {
	cfg := config.GetCachexSettingByName("comment")
	cx, err := infraCachex.NewCacheXBuilderByConfig[int64, *dto.Comment](context.Background(), cfg).
		SetGetRealData((&CommentCachex{}).GetRealDataFn()).
		SetMGetRealData((&CommentCachex{}).MGetRealDataFn()).
		Build()
	if err != nil {
		panic(fmt.Errorf("init cachex error: %w", err))
	}
	commentCachex = &CommentCachex{
		cacheX: cx,
		expire: cfg.GetExpire(),
	}
}

func (c *CommentCachex) GetRealDataFn() func(ctx context.Context, id int64) (*dto.Comment, error) {
	return func(ctx context.Context, id int64) (*dto.Comment, error) {
		comment, err := persistence.SelectApprovalCommentByID(mysql.GetDB(ctx), id)
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		return dto.NewCommentByModel(comment), nil
	}
}

func (c *CommentCachex) MGetRealDataFn() func(ctx context.Context, ids []int64) (map[int64]*dto.Comment, error) {
	return func(ctx context.Context, ids []int64) (map[int64]*dto.Comment, error) {
		commentModels, err := persistence.SelectApprovalCommentsByIDs(mysql.GetDB(ctx), misc.SliceDeduplicate(ids))
		if err != nil {
			return nil, infraCachex.ParseErr(err)
		}
		comments := make(map[int64]*dto.Comment, len(commentModels))
		for _, m := range commentModels {
			comments[m.ID] = dto.NewCommentByModel(m)
		}
		return comments, nil
	}
}

func (c *CommentCachex) MGet(ctx context.Context, ids []int64) map[int64]*dto.Comment {
	return c.cacheX.MGet(ctx, ids, c.expire)
}

func (c *CommentCachex) Get(ctx context.Context, id int64) (*dto.Comment, error) {
	comment, ok := c.cacheX.Get(ctx, id, c.expire)
	if !ok {
		return nil, consts.ErrRecordNotFound
	}
	return comment, nil
}
