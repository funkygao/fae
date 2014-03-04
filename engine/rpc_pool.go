package engine

import (
	log "github.com/funkygao/log4go"
	"sync/atomic"
	"time"
)

type rpcHandler func(req interface{})

// Like php-fpm pm pool
type rpcThreadPool struct {
	cf           *configProcessManagement
	handler      rpcHandler
	spareServerN int32
	reqChan      chan interface{} // max outstanding session throttle
}

func newRpcThreadPool(cf *configProcessManagement,
	handler rpcHandler) (this *rpcThreadPool) {
	this = new(rpcThreadPool)
	this.cf = cf
	this.handler = handler
	this.reqChan = make(chan interface{}, this.cf.maxOutstandingSessions)

	return
}

func (this *rpcThreadPool) start() {
	if this.cf.dynamic() {
		this.spawnChildrenInBatch(this.cf.startServers)
	}
}

func (this *rpcThreadPool) spawnChildrenInBatch(batchSize int) {
	t1 := time.Now()
	for i := 0; i < batchSize; i++ {
		go this.dynamicHandleRequest()
		atomic.AddInt32(&this.spareServerN, 1)
	}

	log.Debug("rpcThreadPool spawned %d children within %s", n, time.Since(t1))
}

func (this *rpcThreadPool) dispatch(request interface{}) {
	if this.cf.dynamic() {
		this.reqChan <- request
	} else {
		// here, reqChan is just a throttle to control max outstanding sessions
		this.reqChan <- true
		go func() {
			this.handler(request)
			<-this.reqChan
		}()
	}
}

func (this *rpcThreadPool) dynamicHandleRequest() {
	for {
		req := <-this.reqChan // will block

		// got a request, before finishing it, I'm not spare
		atomic.AddInt32(&this.spareServerN, -1)

		// spawn children in batch if neccessary
		leftN := atomic.LoadInt32(&this.spareServerN)
		if leftN < this.cf.minSpareServers {
			log.Warn("rpc thread pool seems busy: left %d", leftN)
			go this.spawnChildrenInBatch(this.cf.spawnServers)
		}

		// handle request
		this.handler(req)

		// this request finished, I'm spare again, able to handle new request
		atomic.AddInt32(&this.spareServerN, 1)
	}

}
