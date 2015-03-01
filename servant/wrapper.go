package servant

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	//log "github.com/funkygao/log4go"
)

type FunServantImplWrapper struct {
	*FunServantImpl
}

func NewFunServantWrapper(cf *config.ConfigServant) (this *FunServantImplWrapper) {
	this = &FunServantImplWrapper{FunServantImpl: NewFunServant(cf)}
	return
}

func (this *FunServantImplWrapper) Ping(ctx *rpc.Context) (r string, ex error) {
	r, ex = this.FunServantImpl.Ping(ctx)
	return
}
