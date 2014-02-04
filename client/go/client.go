package main

import (
	"fmt"
	"github.com/funkygao/fae/engine"
	"time"
)

func main() {
	t1 := time.Now()

	client, transport := engine.Client(":9001")
	defer transport.Close()

	for i := 0; i < 10; i++ {
		r, _ := client.Ping()

		fmt.Println(r, time.Since(t1))
		t1 = time.Now()
	}
}
