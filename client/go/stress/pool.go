package main

import (
	"github.com/funkygao/fae/servant/proxy"
	"log"
	"time"
)

func testServantPool(proxy *proxy.Proxy) {
	t1 := time.Now()
	peer := host + ":9001"
	for i := 0; i < Concurrency; i++ {
		client, err := proxy.ServantByAddr(peer)
		if err != nil {
			log.Printf("seq^%d err^%v\n", i, err)
			return
		}

		client.Close()
		client.Recycle()
	}

	log.Printf("proxy.ServantByAddr[%s] open/close %d loops within %s\n",
		peer, Concurrency, time.Since(t1))

	key := "test.go"
	t1 = time.Now()
	for i := 0; i < Concurrency; i++ {
		client, err := proxy.ServantByKey(key)
		if err != nil {
			log.Printf("seq^%d err^%v\n", i, err)
			return
		}

		client.Close()
		client.Recycle()

		peer = client.Addr()
	}

	log.Printf("proxy.ServantByKey[%s] open/close %d loops within %s\n",
		peer, Concurrency, time.Since(t1))
}
