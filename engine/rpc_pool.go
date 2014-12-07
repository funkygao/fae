package engine

type rpcClientHandler func(req interface{})

// Like php-fpm pm pool
// forking goroutine under benchmark is around 40k/s, if higher conn/s
// is required, need pre-fork goroutines
type rpcThreadPool struct {
	handler rpcClientHandler
	reqChan chan interface{} // max outstanding session throttle
}

func newRpcThreadPool(maxOutstandingSessions int,
	handler rpcClientHandler) (this *rpcThreadPool) {
	this = new(rpcThreadPool)
	this.handler = handler
	this.reqChan = make(chan interface{}, maxOutstandingSessions)

	return
}

func (this *rpcThreadPool) Dispatch(request interface{}) {
	this.reqChan <- true // block if outstanding sessions overflows
	go func() {
		this.handler(request)
		<-this.reqChan
	}()
}
