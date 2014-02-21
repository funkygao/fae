/*
dlog ident:string tag:string
*/
package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/dlog"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) Dlog(ctx *rpc.Context, ident string, tag string,
	json string) (appErr error) {
	this.stats.inc("dlog")

	if err := dlog.Dlog(ident, tag, json); err != nil {
		log.Error("dlog: %v", err)
	}

	return nil
}
