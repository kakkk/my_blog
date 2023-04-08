package dto

import "encoding/json"

type PostMeta struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func (a *PostMeta) Serialize() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}

func (a *PostMeta) Deserialize(str string) (*PostMeta, error) {
	meta := &PostMeta{}
	err := json.Unmarshal([]byte(str), meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}
