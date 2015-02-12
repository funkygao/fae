package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/funkygao/log4go"
)

type rpcClientHandler func(req interface{})

var Null struct{}

// Like php-fpm pm pool
// forking goroutine under benchmark is around 40k/s, if higher conn/s
// is required, need pre-fork goroutines
type rpcThreadPool struct {
	preforkMode bool
	handler     rpcClientHandler

	throttleChan     chan struct{}          // if not prefork mode
	clientSocketChan chan thrift.TTransport // if prefork mode
}

func newRpcThreadPool(prefork bool, maxOutstandingSessions int,
	handler rpcClientHandler) (this *rpcThreadPool) {
	this = new(rpcThreadPool)
	this.handler = handler
	this.preforkMode = prefork

	if this.preforkMode {
		this.clientSocketChan = make(chan thrift.TTransport, maxOutstandingSessions)
		for i := 0; i < maxOutstandingSessions; i++ {
			// prefork
			go func() {
				for {
					this.handler(<-this.clientSocketChan)
				}
			}()
		}
	} else {
		this.throttleChan = make(chan struct{}, maxOutstandingSessions)
	}

	return
}

func (this *rpcThreadPool) Dispatch(clientSocket thrift.TTransport) {
	if this.preforkMode {
		select {
		case this.clientSocketChan <- clientSocket:
		default:
			log.Warn("rpc thread pool full, discarded client: %+v", clientSocket)
		}

		return
	}

	this.throttleChan <- Null // block if outstanding sessions overflows
	go func() {
		this.handler(clientSocket)
		<-this.throttleChan
	}()
}
