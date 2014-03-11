package mysql

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
	"sync"
)

type StandardServerSelector struct {
	conf     *config.ConfigMysql
	mu       sync.Mutex
	breakers map[string]*breaker.Consecutive // key is dsn
	clients  map[string]*mysql               // key is dsn
}

func newStandardServerSelector(cf *config.ConfigMysql) (this *StandardServerSelector) {
	this = new(StandardServerSelector)
	this.conf = cf
	this.breakers = make(map[string]*breaker.Consecutive)
	this.clients = make(map[string]*mysql)
	for _, server := range cf.Servers {
		my := newMysql(server.DSN())
		my.setMaxIdleConns(cf.MaxIdleConnsPerServer)
		this.clients[server.DSN()] = my
	}
	return
}

func (this *StandardServerSelector) PickServer(pool string,
	shardId int) (my *mysql, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return
}

func (this *StandardServerSelector) endsWithDigit(pool string) bool {
	lastChar := pool[len(pool)-1]
	return lastChar >= '0' && lastChar <= '9'

}
