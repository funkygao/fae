package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// get a uniq name with length 3
func (this *FunServantImpl) SyncName3(ctx *rpc.Context,
	name string) (r bool, appErr error) {
	const IDENT = "sync.name3"

	this.stats.inc(IDENT)

	if err := this.namegen.SetBusy(name); err != nil {
		r = false

		log.Error("%s[%s]: %s", IDENT, name, err)
	} else {
		r = true
	}

	log.Debug("%s: %s -> %s", IDENT, ctx.Host, name)

	return
}
