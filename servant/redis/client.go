package redis

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/redigo/redis"
	"sync"
	"time"
)

type Client struct {
	breaker *breaker.Consecutive

	selectors map[string]ServerSelector         // key is pool name
	locks     map[string]map[string]*sync.Mutex // pool:serverAddr:Mutex
	conns     map[string]map[string]*redis.Pool // pool:serverAddr:redis.Pool
}

func New(cf *config.ConfigRedis) *Client {
	this := new(Client)
	this.selectors = make(map[string]ServerSelector)
	this.conns = make(map[string]map[string]*redis.Pool)
	this.locks = make(map[string]map[string]*sync.Mutex)
	this.breaker = &breaker.Consecutive{
		FailureAllowance: cf.Breaker.FailureAllowance,
		RetryTimeout:     cf.Breaker.RetryTimeout}
	for pool, _ := range cf.Servers {
		this.selectors[pool] = new(ConsistentServerSelector)
		this.selectors[pool].SetServers(cf.PoolServers(pool)...)
		this.conns[pool] = make(map[string]*redis.Pool)
		this.locks[pool] = make(map[string]*sync.Mutex)
		for _, addr := range cf.PoolServers(pool) {
			this.locks[pool][addr] = &sync.Mutex{}

			this.conns[pool][addr] = &redis.Pool{
				MaxIdle:     cf.Servers[pool][addr].MaxIdle,
				MaxActive:   cf.Servers[pool][addr].MaxActive,
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

func (this *Client) Call(cmd string, pool string, key string, val ...interface{}) (newVal interface{}, err error) {
	addr := this.addr(pool, key)
	conn := this.conns[pool][addr].Get()
	err = conn.Err()
	if err != nil {
		conn.Close()
		this.breaker.Fail()
		log.Error("redis.%s[%s] conn: %s", cmd, key, err)
		return
	}

	this.locks[pool][addr].Lock()

	// Do(cmd string, args ...interface{}) (reply interface{}, err error)
	switch cmd {
	case "GET":
		newVal, err = conn.Do(cmd, key)
		if newVal == nil {
			err = ErrorDataNotExists
		}

	case "SET":
		_, err = conn.Do(cmd, key, val[0])

	case "DEL":
		_, err = conn.Do(cmd, key)
	}
	if err != nil && err != ErrorDataNotExists {
		log.Error("redis.%s[%s]: %s", cmd, key, err)
	}

	this.breaker.Succeed()
	this.locks[pool][addr].Unlock()
	conn.Close() // return to conn pool
	return
}

func (this *Client) Set(pool string, key string, val interface{}) (err error) {
	_, err = this.Call("SET", pool, key, val)
	return
}

// newVal types are represented using the following Go types:
// error                   redis.Error
// integer                 int64
// simple string           string
// bulk string             []byte or nil if value not present.
// array                   []interface{} or nil if value not present
func (this *Client) Get(pool string, key string) (newVal interface{}, err error) {
	newVal, err = this.Call("GET", pool, key)
	return
}

func (this *Client) Del(pool, key string) (err error) {
	_, err = this.Call("DEL", pool, key)
	return
}

func (this *Client) addr(pool, key string) string {
	// FIXME if pool not exists, will panic
	return this.selectors[pool].PickServer(key)
}
