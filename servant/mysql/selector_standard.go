package mysql

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/breaker"
)

type StandardServerSelector struct {
	conf     *config.ConfigMysql
	breakers map[string]*breaker.Consecutive // key is pool
	clients  map[string]*mysql               // key is pool
}

func newStandardServerSelector(cf *config.ConfigMysql) (this *StandardServerSelector) {
	const MAX_RETRIES = 3
	this = new(StandardServerSelector)
	this.conf = cf
	this.breakers = make(map[string]*breaker.Consecutive)
	this.clients = make(map[string]*mysql)
	for _, server := range cf.Servers {
		my := newMysql(server.DSN())
		retries := 0
		for ; retries < MAX_RETRIES; retries++ {
			err := my.Open()
			if err == nil {
				break
			}
		}
		if retries == MAX_RETRIES {
			// FIXME
		}

		my.setMaxIdleConns(cf.MaxIdleConnsPerServer)
		this.clients[server.Pool] = my
		this.breakers[server.Pool] = &breaker.Consecutive{
			FailureAllowance: this.conf.Breaker.FailureAllowance,
			RetryTimeout:     this.conf.Breaker.RetryTimeout}
	}

	return
}

func (this *StandardServerSelector) PickServer(pool string,
	shardId int) (*mysql, error) {
	my, present := this.clients[pool]
	if !present {
		return nil, ErrServerNotFound
	}

	return my, nil
}

func (this *StandardServerSelector) endsWithDigit(pool string) bool {
	lastChar := pool[len(pool)-1]
	return lastChar >= '0' && lastChar <= '9'

}
