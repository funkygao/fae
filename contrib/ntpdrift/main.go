package main

import (
	"log"
	"time"
)

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func main() {
	lastT := milliseconds()
	for {
		ts := milliseconds()
		if ts < lastT {
			log.Printf("NTP drifted clock backwards: %d %d\n", lastT, ts)
		}

		lastT = ts
		time.Sleep(time.Microsecond)
	}

}
