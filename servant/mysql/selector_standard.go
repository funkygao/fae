package mysql

import (
	"github.com/funkygao/fae/config"
)

type StandardServerSelector struct {
	conf    *config.ConfigMysql
	clients map[string]*mysql // key is pool
}

func newStandardServerSelector(cf *config.ConfigMysql) (this *StandardServerSelector) {
	this = new(StandardServerSelector)
	this.conf = cf
	this.clients = make(map[string]*mysql)
	for _, server := range cf.Servers {
		my := newMysql(server.DSN(), &cf.Breaker)
		for retries := uint(0); retries < cf.Breaker.FailureAllowance; retries++ {
			err := my.Open()
			if err == nil {
				break
			}

			my.breaker.Fail()
		}

		my.setMaxIdleConns(cf.MaxIdleConnsPerServer)
		this.clients[server.Pool] = my
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
