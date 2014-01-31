package engine

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	//"github.com/funkygao/fxi/servant"
)

func (this *Engine) ServeForever() {
	listenSocket, err := thrift.NewTServerSocket(this.conf.String("listen_addr", ""))

}
