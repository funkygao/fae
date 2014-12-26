package main

import (
	"flag"
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
)

var (
	host string
	port string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "host name of faed")
	flag.StringVar(&port, "port", "9001", "fae port")
	flag.Parse()
}

func main() {
	client, err := proxy.NewWithDefaultConfig().Servant(host + ":" + port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Recycle()

	ctx := rpc.NewContext()
	ctx.Reason = "pingfae"
	ctx.Rid = "1"
	r, err := client.Ping(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

}
