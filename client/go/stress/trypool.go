package main

import (
	"github.com/funkygao/fae/servant/proxy"
	"log"
	"time"
)

func tryServantPool(proxy *proxy.Proxy) {
	for i := 0; i < Concurrency; i++ {
		t1 := time.Now()
		client, err := proxy.Servant(host + ":9001")
		if err != nil {
			log.Printf("seq^%d err^%v\n", i, err)
			return
		}
		log.Printf("%8d connected within %s", i+1, time.Since(t1))
		client.Recycle()
	}

	log.Println("try server connect/close done!!!")
}
