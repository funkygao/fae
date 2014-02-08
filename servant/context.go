package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"strings"
)

type callerInfo struct {
	httpMethod string
	uri        string
	seqId      string
}

func (this *callerInfo) Valid() bool {
	return this.seqId != ""
}

func (this *FunServantImpl) callerInfo(ctx *rpc.Context) (r callerInfo) {
	const N = 3
	p := strings.SplitN(ctx.Caller, "+", N)
	if len(p) != N {
		return
	}

	r.httpMethod = p[0]
	r.uri = p[1]
	r.seqId = p[2]
	return
}
