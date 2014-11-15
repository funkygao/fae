package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, appErr error) {
	if appErr = validateContext(ctx); appErr != nil {
		return
	}

	this.stats.inc("ping")
	log.Debug("ping from %+v", *ctx)
	return "pong", nil
}
