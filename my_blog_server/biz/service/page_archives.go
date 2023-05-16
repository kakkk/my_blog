package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cast"

	"my_blog/biz/common/config"
	"my_blog/biz/common/errorx"
	"my_blog/biz/common/log"
	"my_blog/biz/common/utils"
	"my_blog/biz/entity"
	"my_blog/biz/model/blog/page"
	"my_blog/biz/repository/mysql"
	"my_blog/biz/repository/storage"
)

func ArchivesPage(ctx context.Context) (rsp *page.ArchivesPageResp, pErr *errorx.PageError) {
	rsp = &page.ArchivesPageResp{
		Meta: &page.PageMeta{
			Title:       fmt.Sprintf("文章归档 - %v", config.GetBlogName()),
			Description: "文章归档",
			CDNDomain:   config.GetSiteConfig().CDNDomain,
			SiteDomain:  config.GetSiteConfig().SiteDomain,
			PageType:    page.PageTypeArchives,
		},
	}
	archives, expired := storage.GetArchivesStorage().Get()
	// 数据过期，异步拉取数据
	if expired {
		go RefreshArchives(ctx)
	}
	rsp.PostArchives = archives
	return rsp, nil
}

func RefreshArchives(ctx context.Context) {
	// 不能使用postMeta，拉取全量数据会破坏LRU
	logger := log.GetLoggerWithCtx(ctx)
	defer utils.Recover(ctx, func() {})()

	postFromDB, err := mysql.SelectAllPublishedPostWithBatch(mysql.GetDB(ctx))
	if err != nil {
		logger.Errorf("select all post error:[%v]", err)
		return
	}
	// map[year][month][]post
	postMap := map[int]map[time.Month][]*entity.Article{}
	// 最早的年份
	nowYear := time.Now().Year()
	minYear := nowYear
	for _, post := range postFromDB {
		pub := post.PublishAt
		if pub.Year() < minYear {
			minYear = pub.Year()
		}
		if _, ok := postMap[pub.Year()]; !ok {
			postMap[pub.Year()] = map[time.Month][]*entity.Article{}
		}
		if _, ok := postMap[pub.Year()][pub.Month()]; !ok {
			postMap[pub.Year()][pub.Month()] = []*entity.Article{}
		}
		postMap[pub.Year()][pub.Month()] = append(postMap[pub.Year()][pub.Month()], post)
	}

	var archives []*page.ArchiveByYear

	for y := nowYear; y >= minYear; y-- {
		if _, ok := postMap[y]; !ok {
			continue
		}
		count := 0
		var byMonth []*page.ArchiveByMonth
		for m := 12; m >= 1; m-- {
			posts, ok := postMap[y][time.Month(m)]
			if !ok {
				continue
			}
			byMonth = append(byMonth, &page.ArchiveByMonth{
				Posts: postListToArchivesPostItem(posts),
				Month: time.Month(m).String(),
				Count: cast.ToString(len(posts)),
			})
			count += len(posts)
		}
		archives = append(archives, &page.ArchiveByYear{
			Archives: byMonth,
			Year:     cast.ToString(y),
			Count:    cast.ToString(count),
		})
	}
	storage.GetArchivesStorage().Set(archives)
	logger.Info("refresh archives success")
}

func postListToArchivesPostItem(posts []*entity.Article) []*page.PostItem {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PublishAt.Unix() >= posts[j].PublishAt.Unix()
	})
	var result []*page.PostItem
	for _, post := range posts {
		result = append(result, &page.PostItem{
			ID:    cast.ToString(post.ID),
			Title: cast.ToString(post.Title),
			Info:  post.PublishAt.Format("January 02, 2006"),
		})
	}
	return result
}
