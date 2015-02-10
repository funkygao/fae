package redis

import (
	"errors"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/redigo/redis"
	"sync"
	"time"
)

// TODO use my own pool, TestOnBorrow is expensive
type Client struct {
	cf      *config.ConfigRedis
	breaker *breaker.Consecutive

	selectors map[string]ServerSelector         // key is pool name
	locks     map[string]map[string]*sync.Mutex // pool:serverAddr:Mutex
	conns     map[string]map[string]*redis.Pool // pool:serverAddr:redis.Pool
}

func New(cf *config.ConfigRedis) *Client {
	this := new(Client)
	this.cf = cf
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

func (this *Client) Call(cmd string, pool string,
	keysAndArgs ...interface{}) (newVal interface{}, err error) {
	if this.breaker.Open() {
		return nil, ErrCircuitOpen
	}

	if len(keysAndArgs) == 0 {
		// e,g. cmd=multi | exec | discard
		return nil, errors.New("redis cmd not implemented:" + cmd)
	}

	key := keysAndArgs[0].(string)
	addr, errPoolNotFound := this.addr(pool, key)
	if errPoolNotFound != nil {
		return nil, errPoolNotFound
	}

	conn := this.conns[pool][addr].Get()
	err = conn.Err()
	if err != nil {
		conn.Close()
		this.breaker.Fail() // conn err is always system err
		return
	}

	this.locks[pool][addr].Lock()

	// Do(cmd string, args ...interface{}) (reply interface{}, err error)
	switch cmd {
	case "GET":
		newVal, err = conn.Do(cmd, key)
		if newVal == nil {
			err = ErrKeyNotExist
		}

	default:
		newVal, err = conn.Do(cmd, keysAndArgs...)
	}
	if err != nil && err != ErrKeyNotExist {
		this.breaker.Fail()
	} else {
		this.breaker.Succeed()
	}

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

func (this *Client) addr(pool, key string) (string, error) {
	selector, present := this.selectors[pool]
	if !present {
		return "", ErrPoolNotFound
	}

	return selector.PickServer(key), nil
}

func (this *Client) Warmup() {
	t1 := time.Now()
	for poolName, pool := range this.conns {
		for addr, conn := range pool {
			log.Debug("redis pool[%s] connecting: %s", poolName, addr)
			for i := 0; i < this.cf.Servers[poolName][addr].MaxActive; i++ {
				c := conn.Get()
				if c.Err() != nil {
					log.Error("redis[%s][%s]: %v", poolName, addr, c.Err())
					continue
				}

				c.Do("PING")
				defer c.Close()
			}
		}
	}
	log.Debug("Redis warmup within %s: %+v",
		time.Since(t1), this.selectors)
}
