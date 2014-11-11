package servant

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

type contextInfo struct {
	ctx *rpc.Context

	httpMethod string
	uri        string
	seqId      string
}

func (this *contextInfo) Valid() bool {
	return this.ctx.Rid != ""
}

func (this contextInfo) String() string {
	if !this.Valid() {
		return "Invalid"
	}

	s := fmt.Sprintf("%s^%s+%s", this.httpMethod, this.uri, this.ctx.Rid)
	if this.ctx.IsSetHost() {
		s = fmt.Sprintf("%s H^%s", s, *this.ctx.Host)
	}
	if this.ctx.IsSetIp() {
		s = fmt.Sprintf("%s I^%s", s, *this.ctx.Ip)
	}
	if this.ctx.IsSetSid() {
		s = fmt.Sprintf("%s S^%s", s, *this.ctx.Sid)
	}
	// currently handles action cmds
	if this.ctx.IsSetReserved() {
		s = fmt.Sprintf("%s A^%s", s, *this.ctx.Reserved)
	}

	return s
}
