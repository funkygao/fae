package servant

import (
	"github.com/funkygao/fae/servant/mongo"
	"github.com/funkygao/golib/trace"
	log "github.com/funkygao/log4go"
)

type mongoProtocolLogger struct {
}

func (this *mongoProtocolLogger) Output(calldepth int, s string) error {
	log.Debug("(%s) %s", trace.CallerFuncName(calldepth), s)
	return nil
}

func (this *FunServantImpl) mongoSession(pool string,
	shardId int32) (*mongo.Session, error) {
	sess, err := this.mg.Session(pool, shardId)
	if err != nil {
		log.Error("{pool^%s id^%d} %s", pool, shardId, err)
		return nil, err
	}

	return sess, err
}
