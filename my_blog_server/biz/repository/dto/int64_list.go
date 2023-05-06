package dto

import "encoding/json"

type Int64List []int64

func (l *Int64List) Serialize() string {
	bytes, _ := json.Marshal(l)
	return string(bytes)
}

func (l *Int64List) Deserialize(str string) (*Int64List, error) {
	list := &Int64List{}
	err := json.Unmarshal([]byte(str), list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (l *Int64List) ToInt64List() []int64 {
	return *l
}

func NewInt64List(list []int64) *Int64List {
	return (*Int64List)(&list)
}
