package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *FunServantImpl) Ping(ctx *rpc.Context) (r string, intError error) {
	log.Debug("ping from %+v", *ctx)
	return "pong", nil
}
