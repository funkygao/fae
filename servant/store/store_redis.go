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
	if err := this.redis.Get(POOL, key, &val); err == nil {
		present = true
	}
	return
}

func (this *RedisStore) Set(key string, val interface{}) {
	this.redis.Set(POOL, key, val)
}

func (this *RedisStore) Del(key string) {
	this.redis.Del(POOL, key)
}
