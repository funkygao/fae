package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/gofmt"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) extractUid(ctx *rpc.Context) (uid int64) {
	if ctx.IsSetUid() {
		uid = *ctx.Uid
	}

	return
}
