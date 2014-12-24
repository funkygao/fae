package store

import (
	"github.com/funkygao/golib/cache"
)

type MemStore struct {
	data *cache.LruCache
}

func NewMemStore(maxEntries int) *MemStore {
	this := &MemStore{data: cache.NewLruCache(maxEntries)}
	return this
}

func (this *MemStore) Get(key string) (interface{}, bool) {
	return this.data.Get(key)
}

func (this *MemStore) Set(key string, val interface{}) {
	this.data.Set(key, val)
}

func (this *MemStore) Del(key string) {
	this.data.Del(key)
}
