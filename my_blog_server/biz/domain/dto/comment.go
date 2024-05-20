package dto

import (
	"fmt"
	"sort"
	"time"

	"my_blog/biz/hertz_gen/blog/api"
	"my_blog/biz/hertz_gen/blog/common"
	"my_blog/biz/infra/config"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/repository/model"
)

type Comment struct {
	ID            int64                `json:"id"`
	ArticleID     int64                `json:"post_id"`
	ReplyID       int64                `json:"reply_id"`
	ParentID      int64                `json:"parent_id"`
	ParentContent string               `json:"parent_content"`
	ArticleTitle  string               `json:"post_title"`
	Nickname      string               `json:"nickname"`
	Email         string               `json:"email"`
	EmailMD5      string               `json:"email_md5"`
	Website       string               `json:"website"`
	Content       string               `json:"content"`
	CreateAt      time.Time            `json:"create_at"`
	Status        common.CommentStatus `json:"status"`
}

func (c *Comment) ToModel() *model.Comment {
	return &model.Comment{
		ID:        c.ID,
		ArticleID: c.ArticleID,
		ReplyID:   c.ReplyID,
		ParentID:  c.ParentID,
		Nickname:  c.Nickname,
		Email:     c.Email,
		Website:   c.Website,
		Content:   c.Content,
		CreateAt:  c.CreateAt,
		Status:    c.Status,
	}
}

func NewCommentByModel(comment *model.Comment) *Comment {
	return &Comment{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		ReplyID:   comment.ReplyID,
		ParentID:  comment.ParentID,
		Nickname:  comment.Nickname,
		Email:     comment.Email,
		EmailMD5:  misc.SumStrMD5(comment.Email),
		Website:   comment.Website,
		Content:   comment.Content,
		CreateAt:  comment.CreateAt,
		Status:    comment.Status,
	}
}

func NewCommentListByModel(comments []*model.Comment) []*Comment {
	result := make([]*Comment, 0, len(comments))
	for _, comment := range comments {
		result = append(result, NewCommentByModel(comment))
	}
	return result
}

func (c *Comment) getCommentAt() string {
	now := time.Now()
	diff := now.Sub(c.CreateAt)
	days := int(diff.Hours() / 24)
	if days >= 1 {
		return fmt.Sprintf("%d days", days)
	}

	hours := int(diff.Hours())
	if hours >= 1 {
		return fmt.Sprintf("%d hours", hours)
	}

	minutes := int(diff.Minutes())
	return fmt.Sprintf("%d minutes", minutes)
}

func (c *Comment) ToResponseEntity() *api.Comment {
	return &api.Comment{
		ID:        c.ID,
		Nickname:  c.Nickname,
		Avatar:    config.GetGravatarCDN() + c.EmailMD5,
		Website:   &c.Website,
		Content:   c.Content,
		CommentAt: c.getCommentAt(),
	}
}

func (c *Comment) ToResponseAdminEntity() *api.GetCommentListAdminItem {
	return &api.GetCommentListAdminItem{
		ID:       c.ID,
		Nickname: c.Nickname,
		Avatar:   config.GetGravatarCDN() + c.EmailMD5,
		Website:  c.Website,
		Article: &api.ArticleMeta{
			ID:    c.ArticleID,
			Title: c.ArticleTitle,
		},
		Content:        c.Content,
		ReplyToID:      &c.ReplyID,
		ReplyToContent: &c.ParentContent,
		Status:         c.Status,
	}
}

func (c *Comment) ToResponseEntityWithAtUser(at string) *api.Comment {
	return &api.Comment{
		ID:        c.ID,
		Nickname:  c.Nickname,
		Avatar:    config.GetGravatarCDN() + c.EmailMD5,
		Website:   &c.Website,
		Content:   c.Content,
		ReplyUser: &at,
		CommentAt: c.getCommentAt(),
	}
}

type Comments []*Comment

func (comments Comments) ToResponseCommentList() []*api.CommentListItem {
	if len(comments) == 0 {
		return make([]*api.CommentListItem, 0)
	}
	var commentList []*api.CommentListItem
	pIDToComments := make(map[int64][]*Comment)
	commentID2Comment := make(map[int64]*Comment)
	// 第一次遍历，找到所有父评论和子评论
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		commentID2Comment[comment.ID] = comment
		if comment.ParentID == 0 {
			item := comment.ToResponseEntity()
			commentList = append(commentList, &api.CommentListItem{Comment: item})
			continue
		}
		_, ok := pIDToComments[comment.ParentID]
		if !ok {
			pIDToComments[comment.ParentID] = []*Comment{comment}
			continue
		}
		pIDToComments[comment.ParentID] = append(pIDToComments[comment.ParentID], comment)
	}
	// 子评论排序
	for pID, c := range pIDToComments {
		sort.Slice(c, func(i, j int) bool {
			return c[i].ID < c[j].ID
		})
		pIDToComments[pID] = c
	}
	// 父评论排序
	sort.Slice(commentList, func(i, j int) bool {
		return commentList[i].Comment.ID < commentList[j].Comment.ID
	})
	// 组装
	for _, item := range commentList {
		list := pIDToComments[item.Comment.ID]
		replies := make([]*api.Comment, 0, len(list))
		for _, comment := range list {
			c, ok := commentID2Comment[comment.ReplyID]
			if !ok {
				continue
			}
			replies = append(replies, comment.ToResponseEntityWithAtUser(c.Nickname))
		}
		item.Replies = replies
	}
	return commentList
}

func (comments Comments) ToResponseAdminCommentList() []*api.GetCommentListAdminItem {
	result := make([]*api.GetCommentListAdminItem, 0, len(comments))
	for _, comment := range comments {
		result = append(result, comment.ToResponseAdminEntity())
	}
	return result
}
