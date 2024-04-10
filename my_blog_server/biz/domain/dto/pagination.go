package dto

import "my_blog/biz/hertz_gen/blog/api"

type Pagination struct {
	HasMore bool
	Total   int64
}

func NewPagination(total int64) *Pagination {
	return &Pagination{Total: total}
}

func (p *Pagination) ToRespPagination(page int32, limit int32) *api.Pagination {
	return &api.Pagination{
		Page:    page,
		Limit:   limit,
		HasMore: p.HasMore,
		Total:   &p.Total,
	}
}
