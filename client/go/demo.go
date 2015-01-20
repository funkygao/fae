package main

import (
	"fmt"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"time"
)

func main() {
	t1 := time.Now()

	client, err := proxy.NewWithDefaultConfig().ServantByAddr(":9001")
	if err != nil {
		panic(err)
	}
	defer client.Recycle()

	ctx := rpc.NewContext()
	ctx.Reason = "gotest"
	ctx.Rid = "189"
	for i := 0; i < 10; i++ {
		r, err := client.Ping(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(r, time.Since(t1))
		t1 = time.Now()
	}

}
