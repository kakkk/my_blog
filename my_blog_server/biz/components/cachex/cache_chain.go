package cachex

import (
	"context"
	"errors"
	"fmt"
)

type cacheChain[T any] struct {
	head *chainNode[T]
	tail *chainNode[T]
}

func newCacheChain[T any]() *cacheChain[T] {
	head := &chainNode[T]{
		cache: &defaultCache[T]{},
	}
	tail := &chainNode[T]{
		cache: &defaultCache[T]{},
	}
	head.next = tail
	tail.prev = head
	chain := &cacheChain[T]{
		head: head,
		tail: tail,
	}
	return chain
}

func (chain *cacheChain[T]) AddCache(_ context.Context, _ string, isSetDefault bool, cache Cache[T]) {
	// 用尾插
	node := &chainNode[T]{
		cache:        cache,
		next:         chain.tail,
		prev:         chain.tail.prev,
		isSetDefault: isSetDefault,
	}
	chain.tail.prev.next = node
	chain.tail.prev = node
}

func (chain *cacheChain[T]) CheckCache(ctx context.Context, name string) error {
	if chain.head.next == chain.tail {
		return errors.New("no cache set")
	}
	curr := chain.head.next
	for curr != chain.tail {
		pong, err := curr.cache.Ping(ctx)
		if err != nil {
			return fmt.Errorf("cache ping error: %w", err)
		}
		logger.Infof(ctx, "cachex [%v] ping [%v]: [%v]", name, curr.cache.Name(), pong)
		curr = curr.next
	}
	return nil
}

func (chain *cacheChain[T]) Get(ctx context.Context, key string) (*CacheData[T], error) {
	return chain.head.next.Get(ctx, key)
}

func (chain *cacheChain[T]) MGet(ctx context.Context, keys []string) (map[string]*CacheData[T], error) {
	return chain.head.next.MGet(ctx, keys)
}

func (chain *cacheChain[T]) Set(ctx context.Context, key string, data *CacheData[T]) error {
	return chain.tail.prev.Set(ctx, key, data)
}

func (chain *cacheChain[T]) MSet(ctx context.Context, kvs map[string]*CacheData[T]) error {
	return chain.tail.prev.MSet(ctx, kvs)
}

func (chain *cacheChain[T]) Delete(ctx context.Context, key string) error {
	return chain.tail.prev.Delete(ctx, key)
}

func (chain *cacheChain[T]) MDelete(ctx context.Context, keys []string) error {
	return chain.tail.prev.MDelete(ctx, keys)
}
