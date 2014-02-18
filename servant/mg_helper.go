package servant

import (
	"github.com/funkygao/fae/servant/mongo"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) mongoSession(pool string,
	shardId int32) (*mongo.Session, error) {
	sess, err := this.mg.Session(pool, shardId)
	if err != nil {
		log.Error("{pool^%s id^%d} %s", pool, shardId, err)
		return nil, err
	}

	return sess, err
}
