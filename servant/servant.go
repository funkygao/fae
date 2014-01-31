/*
Basically it's a plugin framework
*/
package servant

import (
	_ "github.com/funkygao/fxi/servant/memcache" // register memcache pool
	_ "github.com/funkygao/fxi/servant/mongo"    // register mongodb pool
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
