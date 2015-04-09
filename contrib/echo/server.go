// test the NIC PPS rate limit by echo server.
package main

import (
	"fmt"
	"github.com/funkygao/golib/gofmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	CONN_HOST = ""
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var (
	memPool = sync.Pool{New: func() interface{} {
		return make([]byte, 1024)
	}}

	bytesRecved, bytesSent uint64
	clients                int32
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	go showStats()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func showStats() {
	for _ = range time.Tick(time.Second * 10) {
		log.Printf("c:%6d in:%s out:%s",
			atomic.LoadInt32(&clients),
			gofmt.ByteSize(atomic.LoadUint64(&bytesRecved)),
			gofmt.ByteSize(atomic.LoadUint64(&bytesSent)))
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	defer func() {
		conn.Close()
		atomic.AddInt32(&clients, -1)
	}()

	atomic.AddInt32(&clients, 1)

	// Make a buffer to hold incoming data.
	buf := memPool.Get().([]byte)
	for {
		// Read the incoming connection into the buffer.
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading:", err.Error())
			}

			return
		}
		atomic.AddUint64(&bytesRecved, uint64(n))

		// Send a response back to person contacting us.
		n, _ = conn.Write([]byte("Message received."))
		atomic.AddUint64(&bytesSent, uint64(n))
	}

}
