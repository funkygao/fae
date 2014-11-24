package redis

/*
import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/redigo/redis"
	"net"
	"sync"
)

type Client struct {
	conf *config.ConfigMemcache

	selector ServerSelector

	lk       sync.Mutex
	breakers map[net.Addr]*breaker.Consecutive

	pools map[net.Addr]redis.Pool
}

func New(cf config.RedisConfig) *Client {
	this := new(Client)
	this.pool = &redis.Pool{
		MaxIdle:     cf.MaxIdle,
		IdleTimeout: cf.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cf.Server)
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

	return this
}
*/
