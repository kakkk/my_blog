package dto

import "encoding/json"

type StringList []string

func (l *StringList) Serialize() string {
	bytes, _ := json.Marshal(l)
	return string(bytes)
}

func (l *StringList) Deserialize(str string) (*StringList, error) {
	list := &StringList{}
	err := json.Unmarshal([]byte(str), list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (l *StringList) ToStringList() []string {
	return *l
}

func NewStringList(list []string) *StringList {
	return (*StringList)(&list)
}
