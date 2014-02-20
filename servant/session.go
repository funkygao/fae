package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/cache"
)

// most recent session state
type sessions struct {
	items *cache.LruCache
}

func newSessions() *sessions {
	this := &sessions{items: cache.NewLruCache(10 << 10)}
	return this
}

func (this *sessions) set(ctx *rpc.Context, value interface{}) {
	this.items.Set(ctx, value)
}

func (this *sessions) get(ctx *rpc.Context) (interface{}, bool) {
	return this.items.Get(ctx)
}

func (this *sessions) del(ctx *rpc.Context) {
	this.items.Del(ctx)
}
