package memcache

import (
	"github.com/funkygao/fae/config"
	log "github.com/funkygao/log4go"
)

type ClientPool struct {
	conf    *config.ConfigMemcache
	clients map[string]*Client // key is pool name
}

func New(cf *config.ConfigMemcache) *ClientPool {
	this := new(ClientPool)
	this.conf = cf
	this.clients = make(map[string]*Client)
	for _, pool := range cf.Pools() {
		this.clients[pool] = newClient(cf)
	}
	return this
}

func (this *ClientPool) FreeConnMap() map[string]map[string][]*conn {
	ret := make(map[string]map[string][]*conn)
	for pool, client := range this.clients {
		ret[pool] = client.FreeConnMap()
	}
	return ret
}

func (this *ClientPool) Warmup() {
	for _, client := range this.clients {
		client.Warmup()
	}

	log.Debug("Memcache pool warmup finished")
}

func (this *ClientPool) Get(pool string, key string) (item *Item, err error) {
	if client, ok := this.clients[pool]; ok {
		return client.Get(key)
	}
	return nil, ErrInvalidPool
}

func (this *ClientPool) GetMulti(pool string,
	keys []string) (map[string]*Item, error) {
	if client, ok := this.clients[pool]; ok {
		return client.GetMulti(keys)
	}
	return nil, ErrInvalidPool
}

func (this *ClientPool) Set(pool string, item *Item) error {
	if client, ok := this.clients[pool]; ok {
		return client.Set(item)
	}
	return ErrInvalidPool
}

func (this *ClientPool) Add(pool string, item *Item) error {
	if client, ok := this.clients[pool]; ok {
		return client.Add(item)
	}
	return ErrInvalidPool
}

func (this *ClientPool) Increment(pool string, key string,
	delta int64) (newValue uint64, err error) {
	client, ok := this.clients[pool]
	if !ok {
		return 0, ErrInvalidPool
	}

	if delta > 0 {
		return client.Increment(key, uint64(delta))
	}
	return client.Decrement(key, uint64(-delta))
}

func (this *ClientPool) Delete(pool string, key string) error {
	if client, ok := this.clients[pool]; ok {
		return client.Delete(key)
	}
	return ErrInvalidPool
}
