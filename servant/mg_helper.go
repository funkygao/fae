package servant

import (
	log "github.com/funkygao/log4go"
	"strings"
)

func (this *FunServantImpl) mongoSession(kind string,
	shardId int32) (*mongo.Session, error) {
	kind = this.normalizedKind(kind)
	sess, err := this.mg.Session(kind, shardId)
	if err != nil {
		log.Error("{kind^%s id^%d} %s", kind, shardId, err)
		return nil, err
	}

	return sess, err
}

func (this *FunServantImpl) normalizedKind(kind string) string {
	const (
		N      = 2
		PREFIX = "database."
	)

	p := strings.SplitN(kind, PREFIX, N)
	if len(p) == 2 {
		return p[1]
	}

	return kind
}
