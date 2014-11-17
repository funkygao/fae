package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, appErr error) {
	const IDENT = "ping"
	this.stats.inc(IDENT)
	log.Debug("%s from %+v", IDENT, *ctx)
	return "pong", nil
}
