package index

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cast"

	"my_blog/biz/common/utils"
	infraBleve "my_blog/biz/infra/repository/bleve"
	"my_blog/biz/infra/repository/model"
)

type ArticleIndex struct {
	idx bleve.Index
}

type ArticleData struct {
	ID      string
	Title   string
	Content string
}

var articleIndex *ArticleIndex

func GetArticleIndex() *ArticleIndex {
	return articleIndex
}

func InitArticleIndex() error {
	// 当前数据量，内存索引即可
	index, err := infraBleve.NewMemBleveIndex("article.blv")
	if err != nil {
		return fmt.Errorf("new mapping error is: %v", err)
	}
	articleIndex = &ArticleIndex{
		idx: index,
	}
	return nil
}

func (i *ArticleIndex) IndexArticle(article *model.Article) error {
	data := &ArticleData{
		ID:      cast.ToString(article.ID),
		Title:   article.Title,
		Content: utils.StripContent(article.Content),
	}
	err := i.idx.Index(data.ID, data)
	if err != nil {
		return fmt.Errorf("index article error:[%v]", err)
	}
	return nil
}

func (i *ArticleIndex) MIndexArticle(articles []*model.Article) error {
	batch := i.idx.NewBatch()
	for _, article := range articles {
		data := &ArticleData{
			ID:      cast.ToString(article.ID),
			Title:   article.Title,
			Content: utils.StripContent(article.Content),
		}
		err := batch.Index(data.ID, data)
		if err != nil {
			return fmt.Errorf("add batch error:[%v]", err)
		}
	}
	err := i.idx.Batch(batch)
	if err != nil {
		return fmt.Errorf("batch index error:[%v]", err)
	}

	return nil
}

func (i *ArticleIndex) SearchByTitle(queryText string, size int) ([]*ArticleData, error) {
	query := bleve.NewMatchQuery(queryText)
	query.SetField("Title")
	req := bleve.NewSearchRequest(query)
	req.Highlight = bleve.NewHighlight()
	req.Size = size
	resp, err := i.idx.Search(req)
	if err != nil {
		return nil, fmt.Errorf("search error:[%v]", err)
	}
	var result []*ArticleData
	for _, hit := range resp.Hits {
		// 处理标题高亮
		title := ""
		if _, ok := hit.Fragments["Title"]; !ok {
			continue
		}
		for _, hl := range hit.Fragments["Title"] {
			title += hl
		}
		result = append(result, &ArticleData{
			ID:    hit.ID,
			Title: title,
		})
	}
	return result, nil
}
