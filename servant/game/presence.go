package game

import (
	"github.com/funkygao/golib/cache"
	"time"
)

type Presence struct {
	users *cache.LruCache
}

func newPresence() *Presence {
	this := new(Presence)
	this.users = cache.NewLruCache(1 << 20)
	return this
}

func (this *Presence) Update(uid int64) {
	this.users.Set(uid, time.Now().Unix())
}

func (this *Presence) Onlines(uids []int64) []bool {
	const MAX_IDLE = 5 * 60
	r := make([]bool, len(uids))
	now := time.Now().Unix()
	for idx, uid := range uids {
		lastSync, present := this.users.Get(uid)
		if !present {
			r[idx] = false
		} else if now-lastSync.(int64) < MAX_IDLE {
			r[idx] = true
		}
	}
	return r
}
