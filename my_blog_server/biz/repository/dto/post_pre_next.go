package dto

import "encoding/json"

type PostPrevNext struct {
	Prev *int64 `json:"prev"`
	Next *int64 `json:"next"`
}

func (a *PostPrevNext) Serialize() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a *PostPrevNext) Deserialize(str string) (*PostPrevNext, error) {
	pn := &PostPrevNext{}
	err := json.Unmarshal([]byte(str), pn)
	if err != nil {
		return nil, err
	}
	return pn, nil
}
