package mongo

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo"
	"sync"
	"time"
)

type Client struct {
	conf     *config.ConfigMongodb
	selector ServerSelector
	lk       sync.Mutex
	freeconn map[string][]*mgo.Session // the session pool, key is pool

	breaker        *breaker.Consecutive
	connectTimeout time.Duration
	ioTimeout      time.Duration
}

func New(cf *config.ConfigMongodb) (this *Client) {
	this = new(Client)
	this.conf = cf
	this.breaker = &breaker.Consecutive{FailureAllowance: 10,
		RetryTimeout: time.Second * 10}
	this.connectTimeout = time.Duration(this.conf.ConnectTimeout) * time.Second
	this.ioTimeout = time.Duration(this.conf.IoTimeout) * time.Second
	switch cf.ShardStrategy {
	case "legacy":
		this.selector = NewLegacyServerSelector(cf.ShardBaseNum)

	default:
		this.selector = NewStandardServerSelector(cf.ShardBaseNum)
	}
	this.selector.SetServers(cf.Servers)

	go this.runWatchdog()

	return
}

func (this *Client) FreeConn() map[string][]*mgo.Session {
	this.lk.Lock()
	defer this.lk.Unlock()
	return this.freeconn
}

func (this *Client) Session(pool string, shardId int32) (*Session, error) {
	server, err := this.selector.PickServer(pool, int(shardId))
	if err != nil {
		return nil, err
	}

	sess, err := this.getConn(server.Url())
	if err != nil {
		return nil, err
	}

	return &Session{Session: sess, client: this, server: server}, nil
}

func (this *Client) WarmUp() {
	var (
		sess *mgo.Session
		err  error
		t1   = time.Now()
	)
	for retries := 0; retries < 3; retries++ {
		for _, server := range this.selector.ServerList() {
			sess, err = this.getConn(server.Address())
			if err != nil {
				log.Error("Warmup %v fail: %s", server.Address(), err)
				break
			} else {
				this.putFreeConn(server.Url(), sess)
			}
		}

		if err == nil {
			break
		}
	}

	if err == nil {
		log.Trace("Mongodb warmed up within %s: %+v",
			time.Since(t1), this.freeconn)
	} else {
		log.Error("Mongodb failed to warm up within %s: %s",
			time.Since(t1), err)
	}

}

func (this *Client) getConn(url string) (*mgo.Session, error) {
	sess, ok := this.getFreeConn(url)
	if ok {
		return sess, nil
	}

	// create session on demand
	sess, err := mgo.DialWithTimeout(url, this.connectTimeout)
	if err != nil {
		return nil, err
	}

	sess.SetSocketTimeout(this.ioTimeout)
	sess.SetMode(mgo.Monotonic, true)

	return sess, nil
}

func (this *Client) putFreeConn(url string, sess *mgo.Session) {
	this.lk.Lock()
	defer this.lk.Unlock()
	if this.freeconn == nil {
		this.freeconn = make(map[string][]*mgo.Session)
	}
	freelist := this.freeconn[url]
	if len(freelist) >= this.conf.MaxIdleConnsPerServer {
		sess.Close()
		return
	}
	this.freeconn[url] = append(this.freeconn[url], sess)
}

func (this *Client) getFreeConn(url string) (sess *mgo.Session, ok bool) {
	this.lk.Lock()
	defer this.lk.Unlock()
	if this.freeconn == nil {
		return nil, false
	}
	freelist, present := this.freeconn[url]
	if !present || len(freelist) == 0 {
		return nil, false
	}

	// it is no longer free
	sess = freelist[len(freelist)-1] // last item
	this.freeconn[url] = freelist[:len(freelist)-1]
	return sess, true
}

// caller is responsible for lock
func (this *Client) killConn(session *mgo.Session) {
	for addr, sessions := range this.freeconn {
		for idx, sess := range sessions {
			if sess == session { // pointer addr compare
				// https://code.google.com/p/go-wiki/wiki/SliceTricks
				this.freeconn[addr] = append(this.freeconn[addr][:idx],
					this.freeconn[addr][idx+1:]...)
			}
		}
	}
}
