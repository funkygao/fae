package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"strings"
)

type callerInfo struct {
	ctx *rpc.Context

	httpMethod string
	uri        string
	seqId      string
}

func (this *callerInfo) Valid() bool {
	return this.seqId != ""
}

func (this callerInfo) String() string {
	if !this.Valid() {
		return "Invalid Context"
	}

	s := fmt.Sprintf("X{C^%s+%s+%s", this.httpMethod, this.uri, this.seqId)
	if this.ctx.IsSetHost() {
		s = fmt.Sprintf("%s H^%s", s, *this.ctx.Host)
	}
	if this.ctx.IsSetIp() {
		s = fmt.Sprintf("%s I^%s", s, *this.ctx.Ip)
	}
	if this.ctx.IsSetSid() {
		s = fmt.Sprintf("%s S^%s", s, *this.ctx.Sid)
	}

	return s + "}"
}

func (this *FunServantImpl) callerInfo(ctx *rpc.Context) (r callerInfo) {
	const N = 3
	p := strings.SplitN(ctx.Caller, "+", N)
	if len(p) != N {
		return
	}

	r.ctx = ctx
	r.httpMethod, r.uri, r.seqId = p[0], p[1], p[2]

	return
}
