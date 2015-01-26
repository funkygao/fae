package game

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/redis"
	log "github.com/funkygao/log4go"
	_redis "github.com/funkygao/redigo/redis"
)

const (
	REG_USER     = "u"
	REG_KINGDOM  = "k"
	REG_ALLIANCE = "a"
	REG_CHAT     = "c"
)

var (
	regTypes = []string{REG_USER, REG_KINGDOM, REG_ALLIANCE, REG_CHAT}
)

type Register struct {
	cf *config.ConfigGame

	redis *redis.Client

	currentShards map[string]int64 // key is reg type
}

func newRegister(cf *config.ConfigGame, redis *redis.Client) *Register {
	this := new(Register)
	this.cf = cf
	this.currentShards = make(map[string]int64)
	this.redis = redis

	this.loadSnapshot()
	return this
}

func (this *Register) loadSnapshot() {
	const IDENT = "game.reg.loadSnapshot"

	var err error
	for _, typ := range regTypes {
		key := this.currentKey(typ)
		this.currentShards[typ], err = _redis.Int64(this.redis.Call("GET",
			this.cf.RedisServerPool, key))
		if err != nil {
			if err == _redis.ErrNil || err.Error() == redis.ErrKeyNotExist.Error() {
				this.currentShards[typ] = 1
				_, err = this.redis.Call("SET", this.cf.RedisServerPool,
					key, this.currentShards[typ])
				if err != nil {
					log.Error("%s set[%s]: %s", IDENT, key, err.Error())
				}
			} else {
				log.Error("%s: %s", IDENT, err.Error())
			}
		}
	}

	log.Debug("%s done, %+v", IDENT, this.currentShards)
}

func (this *Register) Register(typ string) (int64, error) {
	const IDENT = "register"

	currKey := fmt.Sprintf("reg.%s.%d", typ, this.currentShards[typ])
	n, err := _redis.Int(this.redis.Call("INCR", this.cf.RedisServerPool, currKey))
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

	case REG_CHAT:
		splitThreshold = this.cf.ShardSplit.Chat

	default:
		return 0, ErrInvalidRegType
	}

	if n > splitThreshold {
		this.currentShards[typ]++

		currKey := this.currentKey(typ)
		_, err = this.redis.Call("INCR", this.cf.RedisServerPool, currKey)
		if err != nil {
			log.Error("incr[%s]: %s", currKey, err)
			return 0, err
		}

		key := fmt.Sprintf("reg.%s.%d", typ, this.currentShards[typ])
		_, err = this.redis.Call("INCR", this.cf.RedisServerPool, key)
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
