package main

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"log"
	"time"
)

func main() {

	client, err := proxy.NewWithDefaultConfig().ServantByAddr(":9001")
	if err != nil {
		panic(err)
	}
	defer client.Recycle()

	ctx := rpc.NewContext()
	ctx.Reason = "go.duplex"
	ctx.Rid = "189"
	for i := 0; i < 10; i++ {
		go func() {
			t1 := time.Now()
			r, err := client.Ping(ctx)
			if err != nil {
				log.Println(err)
			} else {
				log.Println(r, time.Since(t1))
			}

		}()
	}

	<-make(chan struct{})
}
