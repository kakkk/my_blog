package dto

type PostPrevNext struct {
	Prev *int64 `json:"prev"`
	Next *int64 `json:"next"`
}
