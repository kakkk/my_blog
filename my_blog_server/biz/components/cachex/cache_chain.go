package cachex

import (
	"context"
	"errors"
	"fmt"
)

type cacheChain struct {
	head *chainNode
	tail *chainNode
}

func newCacheChain() *cacheChain {
	head := &chainNode{}
	trail := &chainNode{}
	chain := &cacheChain{
		head: &chainNode{cache: &defaultCache{}, next: trail},
		tail: &chainNode{cache: &defaultCache{}, prev: head},
	}
	return chain
}

func (chain *cacheChain) AddCache(ctx context.Context, name string, isSetDefault bool, args ...Cache) {
	curr := chain.tail.prev
	for _, cache := range args {
		node := &chainNode{
			cache:        cache,
			next:         curr.next,
			prev:         curr,
			isSetDefault: isSetDefault,
		}
		curr.next = node
		curr = curr.next
		logger.Infof(ctx, "%v add cache: %v, set_default: %v", name, cache.Name(), isSetDefault)
	}
}

func (chain *cacheChain) CheckCache(ctx context.Context, name string) error {
	if chain.head.next == chain.tail {
		return errors.New("no cache set")
	}
	curr := chain.head.next
	for curr != chain.tail {
		pong, err := curr.cache.Ping(ctx)
		if err != nil {
			return fmt.Errorf("cache ping error: %w", err)
		}
		logger.Infof(ctx, "%v ping %v: %v", name, curr.cache.Name(), pong)
	}
	return nil
}

func (chain *cacheChain) Get(ctx context.Context, key string) (*CacheData, error) {
	return chain.head.Get(ctx, key)
}

func (chain *cacheChain) MGet(ctx context.Context, keys []string) (map[string]*CacheData, error) {
	return chain.head.MGet(ctx, keys)
}

func (chain *cacheChain) Set(ctx context.Context, key string, data *CacheData) error {
	return chain.tail.Set(ctx, key, data)
}

func (chain *cacheChain) MSet(ctx context.Context, kvs map[string]*CacheData) error {
	return chain.tail.MSet(ctx, kvs)
}

func (chain *cacheChain) Delete(ctx context.Context, key string) error {
	return chain.tail.Delete(ctx, key)
}

func (chain *cacheChain) MDelete(ctx context.Context, keys []string) error {
	return chain.tail.MDelete(ctx, keys)
}
