package mongo

import (
	"github.com/funkygao/fae/config"
	"labix.org/v2/mgo"
	"sync"
	"time"
)

type Client struct {
	conf    *config.ConfigMongodb
	Timeout time.Duration

	selector ServerSelector

	lk       sync.Mutex
	freeconn map[string][]*mgo.Session
}

func New(cf *config.ConfigMongodb) (this *Client) {
	this = new(Client)
	this.conf = cf
	this.selector = NewStandardServerSelector(cf.ShardBaseNum)
	this.selector.SetServers(cf.Servers)
	return
}
