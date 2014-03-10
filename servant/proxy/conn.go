package proxy

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
)

func connect(serverAddr string) (*rpc.FunServantClient, error) {
	transportFactory := thrift.NewTBufferedTransportFactory(2 << 10)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocketTimeout(serverAddr, 0)
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
