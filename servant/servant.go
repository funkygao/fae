package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fxi/config"
	"github.com/funkygao/fxi/servant/memcache" // register memcache pool
	_ "github.com/funkygao/fxi/servant/mongo"  // register mongodb pool
	"github.com/funkygao/golib/syslogng"
	"time"
)

type FunServantImpl struct {
	conf *config.ConfigServant

	t1 time.Time // timeit

	mc *memcache.Client
}

func NewFunServant(cf *config.ConfigServant) (this *FunServantImpl) {
	this = &FunServantImpl{conf: cf}
	memcacheServers := this.conf.Memcache.ServerList()
	this.mc = memcache.New(this.conf.Memcache.HashStrategy, memcacheServers...)

	log.Debug("memcache servers %v", memcacheServers)
	return
}

func (this *FunServantImpl) Ping() (r string, err error) {
	return "pong", nil
}

func (this *FunServantImpl) McSet(key string, value []byte, expiration int32) (r bool, err error) {
	this.t1 = time.Now()
	err = this.mc.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration})
	if err == nil {
		r = true
	}

	log.Debug("mc_set key:%s value:%s, expiration:%v %s", key, string(value), expiration,
		time.Since(this.t1))

	return
}

func (this *FunServantImpl) McGet(key string) (r []byte, err error) {
	this.t1 = time.Now()
	var it *memcache.Item
	it, err = this.mc.Get(key)
	r = it.Value

	log.Debug("mc_get key:%s %s", key, time.Since(this.t1))
	return
}

func (this *FunServantImpl) Dlog(area string, json string) (err error) {
	syslogng.Printf("%s,%d,%s", area, time.Now().UTC().Unix(), json)
	return nil
}
