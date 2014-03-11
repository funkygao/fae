package mysql

import (
	"github.com/funkygao/fae/config"
)

type StandardServerSelector struct {
}

func newStandardServerSelector() (this *StandardServerSelector) {
	this = new(StandardServerSelector)
	return
}

func (this *StandardServerSelector) SetServers(servers *config.ConfigMysql) {

}

func (this *StandardServerSelector) PickServer(pool string,
	shardId int) (addr string, err error) {
	return
}
