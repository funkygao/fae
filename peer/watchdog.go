package peer

import (
	log "github.com/funkygao/log4go"
	"time"
)

func (this *Peer) runHeartbeat() {
	t := time.NewTicker(time.Duration(this.heartbeatInterval) * time.Second)
	defer t.Stop()

	var msg = peerMessage{}
	var err error
	msg["ip"] = this.selfAddr
	for _ = range t.C {
		if err = this.publish(msg); err != nil {
			log.Error("Publish fails: %v", err)
		}

		for ip, lastAccess := range this.neighbors {
			if time.Since(lastAccess).Seconds() > this.deadThreshold {
				// he is dead
				this.killNeighbor(ip)
			}
		}

	}
}
