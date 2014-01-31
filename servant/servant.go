package servant

import (

//	"github.com/funkygao/fxi/servant/memcache"
//	"github.com/funkygao/fxi/servant/mongo"
)

type FunServantImpl struct {
}

func NewFunServant() (this *FunServantImpl) {
	this = new(FunServantImpl)
	return
}

func (this *FunServantImpl) Ping() (r string, err error) {
	return "pong", nil
}

func init() {

}
