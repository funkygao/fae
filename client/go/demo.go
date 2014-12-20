package main

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/fae/servant/proxy"
	"time"
)

func main() {
	t1 := time.Now()

	cf := config.ConfigProxy{PoolCapacity: 20}
	client, err := proxy.New(cf).Servant(":9001")
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