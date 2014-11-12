package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	log "github.com/funkygao/log4go"
)

// actor lock
// register an uniq player name
// place a new player into a random tile in kingdom map
// under maintenance

func (this *FunServantImpl) gm_maintain(ctx *rpc.Context, pool string) (appErr error) {
	log.Info("mysql maintain: %s", pool)

	return nil
}

func (this *FunServantImpl) gm_register(ctx *rpc.Context, udid string) (appErr error) {
	return nil
}

func (this *FunServantImpl) gm_actor_lockuser(ctx *rpc.Context, uid int64) (appErr error) {
	return nil
}

func (this *FunServantImpl) gm_actor_locktile(ctx *rpc.Context, geohash int64) (appErr error) {
	appErr = ErrNotImplemented
	return nil
}
