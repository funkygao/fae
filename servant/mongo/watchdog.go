package mongo

import (
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo"
	"sync"
	"time"
)

func (this *Client) runWatchdog() {
	ticker := time.NewTicker(time.Duration(this.conf.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	var wg *sync.WaitGroup
	for _ = range ticker.C {
		this.lk.Lock()
		wg = new(sync.WaitGroup)
		for _, sessions := range this.freeconn {
			for _, sess := range sessions {
				wg.Add(1)
				go this.checkServerStatus(wg, sess)
			}
		}
		this.lk.Unlock()

		wg.Wait()
	}

}

// TODO remove dead session from freecon
func (this *Client) checkServerStatus(wg *sync.WaitGroup, sess *mgo.Session) {
	defer wg.Done()
	err := sess.Ping()
	if err != nil {
		log.Error("mongodb err: %v %s", sess.LiveServers(), err)
		sess.Close()
	}
}
