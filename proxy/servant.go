/*
Proxy of remote servant so that we can distribute request
to cluster instead of having to serve all by ourselves.
*/
package proxy

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func (this *Proxy) connect(serverAddr string) (*rpc.FunServantClient, error) {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(serverAddr)
	if err != nil {
		return nil, err
	}

	useTransport := transportFactory.GetTransport(transport)
	client := rpc.NewFunServantClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return nil, err
	}

	return client, nil
}
