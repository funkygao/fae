package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/syslogng"
	"time"
)

func (this *FunServantImpl) Dlog(ctx *rpc.ReqCtx, ident string, tag string,
	json string) (intError error) {
	// add newline and timestamp here
	syslogng.Printf(":%s,%s,%d,%s\n", ident, tag, time.Now().UTC().Unix(), json)

	return nil
}
