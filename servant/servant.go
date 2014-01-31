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

func init() {

}
