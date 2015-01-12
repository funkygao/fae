package game

import (
	"fmt"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/redigo/redis"
	"sync"
	"time"
)

type Register struct {
	redis *redis.Pool
	mutex sync.Mutex

	tiles map[int]map[int]bool

	maxPlayersPerKingdom int
	k                    int // currently all players register to this kingdom
}

func newRegister(maxPlayersPerKingdom int, redisAddr string) *Register {
	this := new(Register)
	this.maxPlayersPerKingdom = maxPlayersPerKingdom
	this.tiles = make(map[int]map[int]bool)
	this.redis = &redis.Pool{
		MaxIdle:     5, // TODO
		MaxActive:   10,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddr)
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

	this.loadSnapshot()
	return this
}

func (this *Register) loadSnapshot() {
	const IDENT = "game.reg.loadSnapshot"

	redisConn := this.redis.Get()
	defer redisConn.Close()

	this.mutex.Lock()
	defer this.mutex.Unlock()

	var err error
	this.k, err = redis.Int(redisConn.Do("GET", "reg.k.curr"))
	if err != nil && err == redis.ErrNil {
		this.k = 1
		redisConn.Do("SET", "reg.k.curr", this.k)
	}

	log.Debug("%s done, k:%d", IDENT, this.k)
}

func (this *Register) RegTile() (k, x, y int) {
	const IDENT = "regtile"

	this.mutex.Lock()
	defer this.mutex.Unlock()

	redisConn := this.redis.Get()
	defer redisConn.Close()

	key := fmt.Sprintf("reg.k.%d", this.k)
	n, err := redis.Int(redisConn.Do("INCR", key))
	if err != nil {
		log.Error("%s: %s", IDENT, err)
		return
	}

	if n > this.maxPlayersPerKingdom {
		this.k++
		redisConn.Do("INCR", "reg.k.curr")

		key = fmt.Sprintf("reg.k.%d", this.k)
		redisConn.Do("INCR", key)

		log.Info("creating new kingdom: %d", this.k)
	}

	k = this.k

	return
}
