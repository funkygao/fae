package servant

import (
	"encoding/json"
	"github.com/funkygao/fae/servant/mongo"
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo/bson"
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

func (this *FunServantImpl) unmarshalBson(d []byte) (v bson.M, err error) {
	err = bson.Unmarshal(d, &v)
	if err != nil {
		log.Error("unmarshalBson: %s", err)
	}

	return
}

func (this *FunServantImpl) unmarshalJson(d []byte) (v bson.M, err error) {
	err = json.Unmarshal(d, &v)
	if err != nil {
		log.Error("unmarshalJson: %s", err)
	}

	return
}
