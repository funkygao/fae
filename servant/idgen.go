package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

// Ticket server
func (this *FunServantImpl) IdNext(ctx *rpc.Context,
	flag int16) (r string, appErr error) {
	r = "hello"
	return
}
