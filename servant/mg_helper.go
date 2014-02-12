package servant

import (
	"github.com/funkygao/fae/servant/mongo"
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) mongoSession(kind string,
	shardId int32) (*mongo.Session, error) {
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		log.Error("{kind^%s id^%d} %s", kind, shardId, err)
		return nil, err
	}

	return sess, err
}
