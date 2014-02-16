package proxy

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"sync"
)

type Proxy struct {
	*sync.RWMutex
	freeServant map[string][]*rpc.FunServantClient
}

func New() (this *Proxy) {
	this = new(Proxy)
	this.RWMutex = new(sync.RWMutex)
	this.freeServant = make(map[string][]*rpc.FunServantClient)
	return
}

func (this *Proxy) Servant(serverAddr string) (*rpc.FunServantClient, error) {
	return this.getServant(serverAddr)
}

func (this *Proxy) getServant(serverAddr string) (*rpc.FunServantClient, error) {
	s, ok := this.getFreeServant(serverAddr)
	if ok {
		return s, nil
	}

	// create the servant connection on demand
	return this.connect(serverAddr)
}

func (this *Proxy) putFreeServant(serverAddr string, client *rpc.FunServantClient) {

}

func (this *Proxy) getFreeServant(serverAddr string) (client *rpc.FunServantClient, ok bool) {
	return
}

// TODO recycle
