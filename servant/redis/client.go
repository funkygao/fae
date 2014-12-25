package redis

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/msgpack"
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

func (this *Client) doCmd(cmd string, pool string, key string, val ...interface{}) (newVal interface{}, err error) {
	addr := this.addr(pool, key)
	conn := this.conns[pool][addr].Get()
	err = conn.Err()
	if err != nil {
		log.Error("redis.%s[%s] conn: %s", cmd, key, err)
		conn.Close()
		this.breaker.Fail()
		return
	}

	this.locks[pool][addr].Lock()

	// Do(cmd string, args ...interface{}) (reply interface{}, err error)
	switch cmd {
	case "GET":
		newVal, err = conn.Do(cmd, key)

	case "SET":
		_, err = conn.Do(cmd, key, val[0])
	}
	if err != nil {
		log.Error("redis.%s[%s]: %s", cmd, key, err)
	}

	this.breaker.Succeed()
	this.locks[pool][addr].Unlock()
	conn.Close() // return to conn pool
	return
}

func (this *Client) Set(pool string, key string, val interface{}) error {
	addr := this.addr(pool, key)
	this.locks[pool][addr].Lock()
	encodedVal, err := msgpack.Marshal(val)
	if err != nil {
		log.Error("msgpack.marshal %+v: %s", val, err)
		this.locks[pool][addr].Unlock()
		return err
	}

	conn := this.conns[pool][addr].Get()
	err = conn.Err()
	if err != nil {
		log.Error("redis.set[%s] socket: %s", key, err)
		this.locks[pool][addr].Unlock()
		conn.Close()
		this.breaker.Fail()
		return err
	}

	// Do(cmd string, args ...interface{}) (reply interface{}, err error)
	if _, err = conn.Do("SET", key, encodedVal); err != nil {
		log.Error("redis.set[%s]: %s", key, err)
	}

	this.breaker.Succeed()
	this.locks[pool][addr].Unlock()
	conn.Close()
	return err
}

func (this *Client) Get(pool string, key string, val interface{}) (err error) {
	addr := this.addr(pool, key)
	this.locks[pool][addr].Lock()

	conn := this.conns[pool][addr].Get()
	err = conn.Err()
	if err != nil {
		log.Error("redis.get[%s] socket: %s", key, err)
		this.locks[pool][addr].Unlock()
		conn.Close()
		this.breaker.Fail()
		return err
	}

	var encodedVal []byte
	encodedVal, err = redis.Bytes(conn.Do("GET", key))
	if encodedVal == nil {
		err = ErrorDataNotExists
		this.locks[pool][addr].Unlock()
		conn.Close()
		return
	}
	if err != nil {
		log.Error("redis.get[%s]: %s", key, err)
		this.locks[pool][addr].Unlock()
		conn.Close()
		return
	}

	this.breaker.Succeed()
	err = msgpack.Unmarshal(encodedVal, val)
	if err != nil {
		log.Error("msgpack.unmarshal %+v: %s", val, err)
	}

	this.locks[pool][addr].Unlock()
	conn.Close()
	return
}

func (this *Client) Del(pool, key string) (err error) {
	addr := this.addr(pool, key)
	this.locks[pool][addr].Lock()

	conn := this.conns[pool][addr].Get()
	err = conn.Err()
	if err != nil {
		log.Error("redis.del[%s] socket: %s", key, err)
		this.locks[pool][addr].Unlock()
		conn.Close()
		this.breaker.Fail()
		return err
	}

	if _, err = conn.Do("DEL", key); err != nil {
		log.Error("redis.del[%s]: %s", key, err)
		this.locks[pool][addr].Unlock()
		conn.Close()
		return err
	}

	this.breaker.Succeed()
	this.locks[pool][addr].Unlock()
	conn.Close()
	return nil
}

func (this *Client) addr(pool, key string) string {
	return this.selectors[pool].PickServer(key)
}
