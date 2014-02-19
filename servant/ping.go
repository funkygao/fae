package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, appErr error) {
	log.Debug("ping from %+v", *ctx)
	this.stats.Ping.Inc(1)
	return "pong", nil
}
