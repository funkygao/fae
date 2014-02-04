package servant

import (
	log "code.google.com/p/log4go"
	"time"
)

func (this *FunServantImpl) runWatchdog() {
	ticker := time.NewTicker(time.Duration(this.conf.WatchdogInterval) * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		log.Debug("lcache len: %d", this.lc.Len())
	}

}
