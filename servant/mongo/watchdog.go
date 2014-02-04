package mongo

import (
	log "code.google.com/p/log4go"
	"labix.org/v2/mgo"
	"sync"
	"time"
)

func (this *Client) runWatchdog() {
	ticker := time.NewTicker(time.Duration(this.conf.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	var wg *sync.WaitGroup
	for _ = range ticker.C {
		log.Debug("mongo servers: %d", len(this.freeconn))

		wg = new(sync.WaitGroup)
		for _, sessions := range this.freeconn {
			for _, sess := range sessions {
				wg.Add(1)
				go this.checkServerStatus(wg, sess)
			}
		}

		wg.Wait()
	}

}

// TODO
func (this *Client) checkServerStatus(wg *sync.WaitGroup, sess *mgo.Session) {
	defer wg.Done()
	err := sess.Ping()
	if err != nil {
		log.Error("mongodb: %s", err)
		sess.Close()
	}
}
