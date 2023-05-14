package cachex

import (
	"context"
)

type CacheData[T any] struct {
	CreateAt      int64 `json:"c"`
	Data          T     `json:"d"`
	IsDefaultData uint  `json:"z"`
}

func (d *CacheData[T]) IsDefault() bool {
	return d.IsDefaultData == 1
}

type Cache[T any] interface {
	Get(ctx context.Context, key string) (*CacheData[T], error)
	MGet(ctx context.Context, keys []string) (map[string]*CacheData[T], error)
	Set(ctx context.Context, key string, data *CacheData[T]) error
	MSet(ctx context.Context, kvs map[string]*CacheData[T]) error
	Delete(ctx context.Context, key string) error
	MDelete(ctx context.Context, keys []string) error
	Ping(ctx context.Context) (string, error)
	Name() string
}

type defaultCache[T any] struct{}

func (*defaultCache[T]) Get(_ context.Context, _ string) (*CacheData[T], error) {
	return nil, ErrNotFound
}
func (*defaultCache[T]) MGet(_ context.Context, _ []string) (map[string]*CacheData[T], error) {
	return map[string]*CacheData[T]{}, ErrNotFound
}
func (*defaultCache[T]) Set(_ context.Context, _ string, _ *CacheData[T]) error {
	return nil
}
func (*defaultCache[T]) MSet(_ context.Context, _ map[string]*CacheData[T]) error {
	return nil
}
func (*defaultCache[T]) Delete(_ context.Context, _ string) error {
	return nil
}
func (*defaultCache[T]) MDelete(_ context.Context, _ []string) error {
	return nil
}

func (*defaultCache[T]) Ping(_ context.Context) (string, error) {
	return "Pong", nil
}
func (c *defaultCache[T]) Name() string {
	return "default"
}
