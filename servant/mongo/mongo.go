package mongo

import (
	"labix.org/v2/mgo"
	"sync"
	"time"
)

type Client struct {
	Timeout time.Duration

	selector ServerSelector

	lk       sync.Mutex
	freeconn map[string][]*mgo.Session
}

func New() (this *Client) {
	this = new(Client)
	return
}
