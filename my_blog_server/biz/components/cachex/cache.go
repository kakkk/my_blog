package cachex

import (
	"context"
)

type CacheData struct {
	CreateAt int64  `json:"c"`
	Data     string `json:"d"`
}

func (d *CacheData) IsDefault() bool {
	return d.Data == ""
}

type Serializable[T any] interface {
	Serialize() string
	Deserialize(string) (T, error)
}

type Cache interface {
	Get(ctx context.Context, key string) (*CacheData, error)
	MGet(ctx context.Context, keys []string) (map[string]*CacheData, error)
	Set(ctx context.Context, key string, data *CacheData) error
	MSet(ctx context.Context, kvs map[string]*CacheData) error
	Delete(ctx context.Context, key string) error
	MDelete(ctx context.Context, keys []string) error
	Ping(ctx context.Context) (string, error)
	Name() string
}

type defaultCache struct{}

func (*defaultCache) Get(ctx context.Context, key string) (*CacheData, error) {
	return nil, ErrNotFound
}
func (*defaultCache) MGet(ctx context.Context, keys []string) (map[string]*CacheData, error) {
	return map[string]*CacheData{}, ErrNotFound
}
func (*defaultCache) Set(ctx context.Context, key string, data *CacheData) error {
	return nil
}
func (*defaultCache) MSet(ctx context.Context, kvs map[string]*CacheData) error {
	return nil
}
func (*defaultCache) Delete(ctx context.Context, key string) error {
	return nil
}
func (*defaultCache) MDelete(ctx context.Context, keys []string) error {
	return nil
}

func (*defaultCache) Ping(ctx context.Context) (string, error) {
	return "Pong", nil
}
func (*defaultCache) Name() string {
	return "default"
}
