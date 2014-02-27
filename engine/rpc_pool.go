package engine

import (
	log "github.com/funkygao/log4go"
	"github.com/funkygao/metrics"
	"time"
)

type rpcHandler func(req interface{})

type rpcThreadPool struct {
	cf           *configProcessManagement
	handler      rpcHandler
	reqChan      chan interface{}
	spareServerN metrics.Counter
}

func newRpcThreadPool(cf *configProcessManagement,
	handler rpcHandler) (this *rpcThreadPool) {
	this = new(rpcThreadPool)
	this.cf = cf
	this.handler = handler
	this.reqChan = make(chan interface{}, this.cf.maxOutstandingSessions)

	// stats
	this.spareServerN = metrics.NewCounter()
	metrics.Register("pool.spare_server", this.spareServerN)
	return
}

func (this *rpcThreadPool) start() {
	if this.cf.dynamic() {
		this.spawnChildren(this.cf.startServers)
	}
}

func (this *rpcThreadPool) spawnChildren(n int) {
	t1 := time.Now()
	for i := 0; i < n; i++ {
		go this.handleRequest()
		this.spareServerN.Inc(1)
	}

	log.Debug("rpcThreadPool spawned %d children within %s", n, time.Since(t1))
}

func (this *rpcThreadPool) dispatch(request interface{}) {
	if this.cf.dynamic() {
		this.reqChan <- request
	} else {
		go this.handler(request)
	}
}

func (this *rpcThreadPool) handleRequest() {
	for {
		req := <-this.reqChan // will block

		// maintain pool spare servers
		this.spareServerN.Dec(1)
		leftN := this.spareServerN.Count()
		if leftN < this.cf.minSpareServers {
			log.Warn("rpc thread pool seems busy: left %d", leftN)
			go this.spawnChildren(this.cf.spawnServers)
		}

		// handle request
		this.handler(req)

		this.spareServerN.Inc(1)
	}

}
