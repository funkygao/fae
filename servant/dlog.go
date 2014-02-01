package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/golib/syslogng"
	"time"
)

func (this *FunServantImpl) Dlog(area string, json string) (err error) {
	this.t1 = time.Now()
	syslogng.Printf("%s,%d,%s", area, time.Now().UTC().Unix(), json)

	log.Debug("dlog area:%s %s", area, time.Since(this.t1))
	return nil
}
