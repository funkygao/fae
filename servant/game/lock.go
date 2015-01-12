package game

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/cache"
	log "github.com/funkygao/log4go"
	"sync"
	"time"
)

type Lock struct {
	cf *config.ConfigLock

	items *cache.LruCache // key: mtime
	mutex sync.Mutex      // lru get/set is safe, but we need more lock span
}

func newLock(cf *config.ConfigLock) *Lock {
	this := &Lock{cf: cf}
	this.items = cache.NewLruCache(cf.MaxItems)
	return this
}

func (this *Lock) Lock(key string) (success bool) {
	this.mutex.Lock()

	mtime, present := this.items.Get(key)
	if !present {
		this.items.Set(key, time.Now())

		this.mutex.Unlock()
		return true
	}

	// present, check expires
	elapsed := time.Since(mtime.(time.Time))
	if this.cf.Expires > 0 && elapsed > this.cf.Expires {
		log.Warn("lock[%s] expires: %s, kicked", key, elapsed)

		// ignore the aged lock, refresh the lock
		this.items.Set(key, time.Now())

		this.mutex.Unlock()
		return true
	}

	this.mutex.Unlock()
	return false
}

func (this *Lock) Unlock(key string) {
	this.items.Del(key)
}
