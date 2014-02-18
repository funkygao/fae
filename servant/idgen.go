package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

// Ticket server
func (this *FunServantImpl) IdNext(ctx *rpc.Context,
	flag int16) (r int64, backwards *rpc.TIdTimeBackwards, appErr error) {
	r, appErr = this.idgen.Next()
	if appErr != nil {
		backwards = appErr.(*rpc.TIdTimeBackwards)
		appErr = nil
	}

	return
}
