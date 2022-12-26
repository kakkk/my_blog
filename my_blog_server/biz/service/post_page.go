package service

import (
	"context"
	"net/http"

	"my_blog/biz/common/config"
	"my_blog/biz/common/consts"
	"my_blog/biz/common/log"
	"my_blog/biz/common/resp"
	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/mysql"

	"github.com/spf13/cast"
)

func PostPage(ctx context.Context, req *page.PostPageRequest) (int, *page.PostPageResponse) {
	logger := log.GetLoggerWithCtx(ctx).WithField("post_id", req.GetID())
	db := mysql.GetDB(ctx)
	// 获取post
	post, err := mysql.SelectPostWithPublishByID(db, req.GetID())
	if err != nil {
		if err == consts.ErrRecordNotFound {
			logger.Warnf("post not found")
			return http.StatusNotFound, &page.PostPageResponse{
				Meta: resp.NewNotFoundErrorMeta(),
			}
		}
		logger.Errorf("select post error:[%v]", err)
		return http.StatusNotFound, &page.PostPageResponse{
			Meta: resp.NewInternalErrorMeta(),
		}
	}

	prevPost, nextPost, err := mysql.SelectPrevNextPostByPublishAt(db, *post.PublishAt)
	if err != nil {
		logger.Errorf("select prev next post error:[%v]", err)
	}

	editor := config.DefaultUserName
	user, err := mysql.SelectUserByID(db, post.CreateUser)
	if err != nil {
		logger.Errorf("select user error:[%v]", err)
	}
	if user != nil {
		editor = user.Nickname
	}

	// 获取标签
	tags, err := getTagListByArticleID(ctx, req.GetID())
	if err != nil {
		logger.Errorf("select tags error:[%v]", err)
	}

	return http.StatusOK, &page.PostPageResponse{
		Title:    post.Title,
		Info:     utils.GetPostInfo(editor, *post.PublishAt, post.Content, post.PV),
		Content:  utils.RenderContent(post.Content),
		Tags:     tags,
		PrevPage: convertPostToPostNav(prevPost),
		NextPage: convertPostToPostNav(nextPost),
		Meta: resp.NewSuccessPageMeta(
			utils.GetPostPageTitle(post.Title),
			utils.GetPostPageDescription(post.Content),
			page.PageTypePost,
		),
	}

}

func convertPostToPostNav(post *entity.Article) *page.PostNav {
	if post == nil {
		return &page.PostNav{
			ID:    "0",
			Title: "",
		}
	}
	return &page.PostNav{
		ID:    cast.ToString(post.ID),
		Title: post.Title,
	}
}
