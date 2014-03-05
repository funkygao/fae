package engine

import (
	log "github.com/funkygao/log4go"
	"sync/atomic"
	"time"
)

type rpcClientHandler func(req interface{})

// Like php-fpm pm pool
// goroutine under benchmark is around 40k/s, if higher conn/s
// is required, need pre-fork goroutines
type rpcThreadPool struct {
	cf           *configProcessManagement
	handler      rpcClientHandler
	spareServerN int32
	reqChan      chan interface{} // max outstanding session throttle
}

func newRpcThreadPool(cf *configProcessManagement,
	handler rpcClientHandler) (this *rpcThreadPool) {
	this = new(rpcThreadPool)
	this.cf = cf
	this.handler = handler
	this.reqChan = make(chan interface{}, this.cf.maxOutstandingSessions)

	return
}

func (this *rpcThreadPool) Start() {
	if this.cf.dynamic() {
		this.spawnDynamicChildrenInBatch(this.cf.startServers)
	}
}

func (this *rpcThreadPool) Dispatch(request interface{}) {
	if this.cf.dynamic() {
		this.reqChan <- request
	} else {
		// here, reqChan is just a throttle to control max outstanding sessions
		this.reqChan <- true // block if outstanding sessions overflows
		go func() {
			this.handler(request)
			<-this.reqChan
		}()
	}
}

func (this *rpcThreadPool) spawnDynamicChildrenInBatch(batchSize int) {
	t1 := time.Now()
	for i := 0; i < batchSize; i++ {
		go this.dynamicHandleRequest()
		atomic.AddInt32(&this.spareServerN, 1)
	}

	log.Debug("rpcThreadPool spawned %d children within %s", batchSize, time.Since(t1))
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
			go this.spawnDynamicChildrenInBatch(this.cf.spawnServers)
		}

		// handle client request
		this.handler(req)

		// this request finished, I'm spare again, able to handle new request
		atomic.AddInt32(&this.spareServerN, 1)
	}

}
