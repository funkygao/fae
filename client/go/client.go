package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"time"
)

func main() {
	t1 := time.Now()

	client, transport := getClient(":9001")
	defer transport.Close()

	for i := 0; i < 10; i++ {
		r, _ := client.Ping()

		fmt.Println(r, time.Since(t1))
		t1 = time.Now()
	}

}

func getClient(serverAddr string) (*rpc.FunServantClient, thrift.TTransport) {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(serverAddr)
	if err != nil {
		panic(err)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := rpc.NewFunServantClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		panic(err)
	}

	return client, transport
}
