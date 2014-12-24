package store

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/redis"
)

const (
	POOL = "default"
)

type RedisStore struct {
	redis *redis.Client
}

func NewRedisStore(cf *config.ConfigRedis) *RedisStore {
	this := &RedisStore{redis: redis.New(cf)}
	return this
}

func (this *RedisStore) Get(key string) (val interface{}, present bool) {
	err := this.redis.Get(POOL, key, &val)
	if err == redis.ErrorDataNotExists {
		present = false
	}
	return
}

func (this *RedisStore) Put(key string, val interface{}) {
	this.redis.Set(POOL, key, val)
}

func (this *RedisStore) Del(key string) {
	this.redis.Del(POOL, key)
}
