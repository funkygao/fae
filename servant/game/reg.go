package game

import (
	"fmt"
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
	"github.com/funkygao/redigo/redis"
	"sync"
	"time"
)

const (
	REG_USER     = "u"
	REG_KINGDOM  = "k"
	REG_ALLIANCE = "a"
)

var (
	regTypes = []string{REG_USER, REG_KINGDOM, REG_ALLIANCE}
)

type Register struct {
	cf *config.ConfigGame

	redis *redis.Pool
	mutex sync.Mutex

	currentShards map[string]int // key is user|alliance|kingdom
}

func newRegister(cf *config.ConfigGame) *Register {
	this := new(Register)
	this.cf = cf
	this.currentShards = make(map[string]int)
	this.redis = &redis.Pool{
		MaxIdle:     5, // TODO
		MaxActive:   10,
		IdleTimeout: 0,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", this.cf.RedisServerAddr)
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
	for _, typ := range regTypes {
		key := this.currentKey(typ)
		this.currentShards[typ], err = redis.Int(redisConn.Do("GET", key))
		if err != nil && err == redis.ErrNil {
			this.currentShards[typ] = 1
			redisConn.Do("SET", key, this.currentShards[typ]) // TODO check err
		}
	}

	log.Debug("%s done, %+v", IDENT, this.currentShards)
}

func (this *Register) Register(typ string) (int, error) {
	const IDENT = "register"

	this.mutex.Lock()
	defer this.mutex.Unlock()

	redisConn := this.redis.Get()
	defer redisConn.Close()

	currKey := fmt.Sprintf("reg.%s.%d", typ, this.currentShards[typ])
	n, err := redis.Int(redisConn.Do("INCR", currKey))
	if err != nil {
		log.Error("%s[%s]: %s", IDENT, typ, err)
		return 0, err
	}

	var splitThreshold int
	switch typ {
	case REG_USER:
		splitThreshold = this.cf.ShardSplit.User

	case REG_ALLIANCE:
		splitThreshold = this.cf.ShardSplit.Alliance

	case REG_KINGDOM:
		splitThreshold = this.cf.ShardSplit.Kingdom

	default:
		return 0, ErrInvalidRegType
	}

	if n > splitThreshold {
		this.currentShards[typ]++

		currKey := this.currentKey(typ)
		_, err = redisConn.Do("INCR", currKey)
		if err != nil {
			log.Error("incr[%s]: %s", currKey, err)
			return 0, err
		}

		key := fmt.Sprintf("reg.%s.%d", typ, this.currentShards[typ])
		_, err = redisConn.Do("INCR", key)
		if err != nil {
			log.Error("incr[%s]: %s", key, err)
			return 0, err
		}
	}

	return this.currentShards[typ], nil
}

func (this *Register) currentKey(typ string) string {
	return fmt.Sprintf("reg.%s.curr", typ)
}
