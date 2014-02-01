package servant

import (
	log "code.google.com/p/log4go"
	"github.com/funkygao/fxi/config"
	"github.com/funkygao/fxi/servant/memcache" // register memcache pool
	_ "github.com/funkygao/fxi/servant/mongo"  // register mongodb pool
)

type FunServantImpl struct {
	conf *config.ConfigServant

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
