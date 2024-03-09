package bleve

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	gseb "github.com/vcaesar/gse-bleve"
)

// NewMemBleveIndex 创建一个Bleve内存索引
func NewMemBleveIndex(indexName string) (bleve.Index, error) {
	opt := gseb.Option{
		Index: indexName,
		Dicts: "embed, zh",
		Opt:   "search-hmm",
		Trim:  "trim",
	}
	index, err := gseb.NewMem(opt)
	if err != nil {
		return nil, fmt.Errorf("new mapping error is: %v", err)
	}
	return index, nil
}
