package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/syslogng"
	"time"
)

func (this *FunServantImpl) Dlog(ctx *rpc.ReqCtx, ident string, tag string,
	json string) (err error) {
	this.t1 = time.Now()

	// add newline and timestamp here
	syslogng.Printf(":%s,%s,%d,%s\n", ident, tag, time.Now().UTC().Unix(), json)

	log.Debug("dlog tag:%s %s", tag, time.Since(this.t1))
	return nil
}
