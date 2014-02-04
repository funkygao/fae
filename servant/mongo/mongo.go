package mongo

import (
	"github.com/funkygao/fae/config"
	"labix.org/v2/mgo"
	"sync"
	"time"
)

type Client struct {
	conf     *config.ConfigMongodb
	selector ServerSelector
	lk       sync.Mutex
	freeconn map[string][]*mgo.Session // the session pool, key is kind

	connectTimeout time.Duration
	ioTimeout      time.Duration
}

func New(cf *config.ConfigMongodb) (this *Client) {
	this = new(Client)
	this.conf = cf
	this.connectTimeout = time.Duration(this.conf.ConnectTimeout) * time.Second
	this.ioTimeout = time.Duration(this.conf.IoTimeout) * time.Second
	this.selector = NewStandardServerSelector(cf.ShardBaseNum)
	this.selector.SetServers(cf.Servers)

	go this.runWatchdog()

	return
}

func (this *Client) Session(kind string, shardId int32) (*Session, error) {
	addr, err := this.selector.PickServer(kind, int(shardId))
	if err != nil {
		return nil, err
	}

	sess, err := this.getConn(addr)
	if err != nil {
		return nil, err
	}

	return &Session{Session: sess, client: this, addr: addr}, nil
}

func (this *Client) getConn(addr string) (*mgo.Session, error) {
	sess, ok := this.getFreeConn(addr)
	if ok {
		return sess, nil
	}

	// create session on demand
	sess, err := mgo.DialWithTimeout(addr, this.connectTimeout)
	if err != nil {
		return nil, err
	}

	sess.SetSocketTimeout(this.ioTimeout)

	return sess, nil
}

func (this *Client) putFreeConn(addr string, sess *mgo.Session) {
	this.lk.Lock()
	defer this.lk.Unlock()
	if this.freeconn == nil {
		this.freeconn = make(map[string][]*mgo.Session)
	}
	freelist := this.freeconn[addr]
	if len(freelist) >= this.conf.MaxIdleConnsPerServer {
		sess.Close()
		return
	}
	this.freeconn[addr] = append(this.freeconn[addr], sess)
}

func (this *Client) getFreeConn(addr string) (sess *mgo.Session, ok bool) {
	this.lk.Lock()
	defer this.lk.Unlock()
	if this.freeconn == nil {
		return nil, false
	}
	freelist, present := this.freeconn[addr]
	if !present || len(freelist) == 0 {
		return nil, false
	}

	// it is no longer free
	sess = freelist[len(freelist)-1] // last item
	this.freeconn[addr] = freelist[:len(freelist)-1]
	return sess, true
}
