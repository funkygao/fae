/*
dlog ident:string tag:string json:json string
*/
package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/syslogng"
	log "github.com/funkygao/log4go"
	"time"
)

func (this *FunServantImpl) Dlog(ctx *rpc.Context, ident string, tag string,
	json string) (intError error) {
	// add newline and timestamp here
	// because signature is void, intError is ignored by client
	if _, intError = syslogng.Printf(":%s,%s,%d,%s\n", ident, tag,
		time.Now().UTC().Unix(), json); intError != nil {
		log.Error(intError)
	}

	return nil
}
