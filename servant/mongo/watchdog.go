package mongo

import (
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo"
	"sync"
	"time"
)

func (this *Client) runWatchdog() {
	if this.conf.HeartbeatInterval == 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(this.conf.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	var wg *sync.WaitGroup
	for _ = range ticker.C {
		wg = new(sync.WaitGroup)
		this.lk.Lock()
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

func (this *Client) checkServerStatus(wg *sync.WaitGroup, sess *mgo.Session) {
	defer wg.Done()
	err := sess.Ping()
	if err != nil {
		// TODO show mongodb url in log
		log.Error("mongodb killed for: %s", err)

		sess.Close()
		this.killConn(sess)
	}
}
