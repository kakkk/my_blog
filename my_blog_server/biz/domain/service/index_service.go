package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cast"
	stripmd "github.com/writeas/go-strip-markdown"

	"my_blog/biz/domain/entity"
	"my_blog/biz/domain/repo/cache"
	"my_blog/biz/domain/repo/persistence"
	"my_blog/biz/hertz_gen/blog/page"
	"my_blog/biz/infra/misc"
	"my_blog/biz/infra/pkg/log"
	repoBleve "my_blog/biz/infra/repository/bleve"
	"my_blog/biz/infra/repository/model"
	"my_blog/biz/infra/repository/mysql"
)

type IndexService interface {
	IndexArticle(article *entity.ArticleData) error
	DeleteArticleIndex(articleID int64) error
	SearchArticleByTitle(queryText string, size int) ([]*entity.ArticleData, error)
	RefreshArchives(ctx context.Context, articles []*model.Article)
	GetArchives(ctx context.Context) []*page.ArchiveByYear
	InitIndex(ctx context.Context) error
}

type IndexServiceImpl struct{}

func GetIndexService() IndexService {
	return &IndexServiceImpl{}
}

// IndexArticle 索引文章
func (svc *IndexServiceImpl) IndexArticle(article *entity.ArticleData) error {
	data := &entity.ArticleIndexData{
		ID:      cast.ToString(article.ID),
		Title:   article.Title,
		Content: svc.stripContent(article.Content),
	}
	err := repoBleve.GetArticleIndex().Index(data.ID, data)
	if err != nil {
		return fmt.Errorf("index article error:[%v]", err)
	}
	return nil
}

// IndexAllPost 批量索引文章
func (svc *IndexServiceImpl) mIndexArticle(articles []*model.Article) error {
	batch := repoBleve.GetArticleIndex().NewBatch()
	for _, article := range articles {
		data := &entity.ArticleIndexData{
			ID:      cast.ToString(article.ID),
			Title:   article.Title,
			Content: misc.StripContent(article.Content),
		}
		err := batch.Index(data.ID, data)
		if err != nil {
			return fmt.Errorf("add batch error:[%v]", err)
		}
	}
	err := repoBleve.GetArticleIndex().Batch(batch)
	if err != nil {
		return fmt.Errorf("batch index error:[%v]", err)
	}

	return nil
}

func (svc *IndexServiceImpl) DeleteArticleIndex(articleID int64) error {
	err := repoBleve.GetArticleIndex().Delete(cast.ToString(articleID))
	if err != nil {
		return fmt.Errorf("delete index fail: %v", err)
	}
	return nil
}

func (svc *IndexServiceImpl) SearchArticleByTitle(queryText string, size int) ([]*entity.ArticleData, error) {
	query := bleve.NewMatchQuery(queryText)
	query.SetField("Title")
	req := bleve.NewSearchRequest(query)
	req.Highlight = bleve.NewHighlight()
	req.Size = size
	resp, err := repoBleve.GetArticleIndex().Search(req)
	if err != nil {
		return nil, fmt.Errorf("search error:[%v]", err)
	}
	var result []*entity.ArticleData
	for _, hit := range resp.Hits {
		// 处理标题高亮
		title := ""
		if _, ok := hit.Fragments["Title"]; !ok {
			continue
		}
		for _, hl := range hit.Fragments["Title"] {
			title += hl
		}
		result = append(result, &entity.ArticleData{
			ID:    cast.ToInt64(hit.ID),
			Title: title,
		})
	}
	return result, nil
}

func (svc *IndexServiceImpl) GetArchives(ctx context.Context) []*page.ArchiveByYear {
	archives, expired := cache.GetArchivesStorage().Get()
	// 数据过期，异步拉取数据
	if expired {
		go svc.RefreshArchives(ctx, nil)
	}
	return archives
}

// 刷新文章归档数据
func (svc *IndexServiceImpl) RefreshArchives(ctx context.Context, articles []*model.Article) {
	logger := log.GetLoggerWithCtx(ctx)
	defer misc.Recover(ctx, func() {})()

	// 未传参，拉取全量数据
	// 此处不能使用postMeta，拉取全量数据会破坏LRU
	if len(articles) == 0 {
		var err error
		articles, err = persistence.SelectAllPublishedPostWithBatch(mysql.GetDB(ctx))
		if err != nil {
			logger.Errorf("select all post error:[%v]", err)
			return
		}
	}
	// map[year][month][]post
	postMap := map[int]map[time.Month][]*model.Article{}
	// 最早的年份
	nowYear := time.Now().Year()
	minYear := nowYear
	for _, post := range articles {
		pub := post.PublishAt
		if pub.Year() < minYear {
			minYear = pub.Year()
		}
		if _, ok := postMap[pub.Year()]; !ok {
			postMap[pub.Year()] = map[time.Month][]*model.Article{}
		}
		if _, ok := postMap[pub.Year()][pub.Month()]; !ok {
			postMap[pub.Year()][pub.Month()] = []*model.Article{}
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
				Posts: svc.postListToArchivesPostItem(posts),
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
	cache.GetArchivesStorage().Set(archives)
	logger.Info("refresh archives success")
}

// InitIndex 初始化索引，用于服务启动时初始化
func (svc *IndexServiceImpl) InitIndex(ctx context.Context) error {
	logger := log.GetLoggerWithCtx(ctx)

	// 加载全量文章数据
	postFromDB, err := persistence.SelectAllPublishedPostWithBatch(mysql.GetDB(ctx))
	if err != nil {
		return fmt.Errorf("load all post from db error:[%v]", err)
	}
	logger.Info("load all post from db success")

	// 索引文章
	err = svc.mIndexArticle(postFromDB)
	if err != nil {
		return fmt.Errorf("index fail: %v", err)
	}
	logger.Infof("index all post success")

	// 刷新归档
	svc.RefreshArchives(ctx, postFromDB)

	return nil
}

func (svc *IndexServiceImpl) stripContent(content string) string {
	striped := stripmd.Strip(content)
	striped = strings.Replace(striped, "\n", " ", -1)
	return striped
}

func (svc *IndexServiceImpl) postListToArchivesPostItem(posts []*model.Article) []*page.PostItem {
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
