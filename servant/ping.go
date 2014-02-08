package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, intError error) {
	log.Debug("ping from %+v", *ctx)
	return "pong", nil
}
