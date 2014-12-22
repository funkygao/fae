package redis

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/msgpack"
	"github.com/funkygao/redigo/redis"
	"net"
	"sync"
)

type Client struct {
	breaker *breaker.Consecutive

	selectors map[string]ServerSelector      // key is pool name
	locks     map[string][string]*sync.Mutex // pool:serverAddr:Mutex
	conns     map[string][string]redis.Pool  // pool:serverAddr:redis.Pool
}

func New(cf config.ConfigRedis) *Client {
	this := new(Client)
	this.selectors = make(map[string]ServerSelector)
	this.conf = make(map[string][string]redis.Pool)
	this.locks = make(map[string][string]*sync.Mutex)
	this.breaker = &breaker.Consecutive{
		FailureAllowance: cf.Breaker.FailureAllowance,
		RetryTimeout:     cf.Breaker.RetryTimeout}
	for pool, _ := range this.conf.Servers {
		this.selectors[pool] = new(ConsistentServerSelector)
		this.selectors[pool].SetServers(cf.PoolServers(pool))
		this.conns[pool] = make(map[string]redis.Pool)
		this.locks[pool] = make(map[string]*sync.Mutex)
		for _, addr := range cf.PoolServers(pool) {
			this.locks[pool][addr] = &sync.Mutex{}

			this.conns[pool][addr] = &redis.Pool{
				MaxIdle:     cf.Servers[pool][addr].MaxIdle,
				IdleTimeout: cf.Servers[pool][addr].IdleTimeout,
				Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", addr)
					if err != nil {
						return nil, err
					}

					return c, err
				},
				TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
				},
			}
		}
	}

	return this
}

func (this *Client) Set(pool string, key string, val interface{}) {
	addr := this.addr(pool, key)
	this.locks[pool][addr].Lock()

	m, err := msgpack.Marshal(val)
	if err != nil {
		log.Error(err)
	}
	_, err = this.conns[pool][addr].Get().Do("SET", key, m)
	if err != nil {
		log.Error(err)
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}

	this.locks[pool][addr].Unlock()
}

func (this *Client) Get(pool string, key string, val interface{}) {
	addr := this.addr(pool, key)
	this.locks[pool][addr].Lock()

	m, err := redis.Bytes(this.conns[pool][addr].Get().Do("GET", key))
	if err != nil {
		log.Error(err)
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}
	err = msgpack.Unmarshal(m, val)
	if err != nil {
		log.Error(err)
	}

	this.locks[pool][addr].Unlock()
}

func (this *Client) Del(pool, key string) {
	addr := this.addr(pool, key)
	this.locks[pool][addr].Lock()

	if _, err := this.conns[pool][addr].Get().Do("DEL", key); err != nil {
		log.Error(err)
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}

	this.locks[pool][addr].Unlock()
}

func (this *Client) addr(pool, key string) string {
	return this.selectors[pool].PickServer(key)
}
