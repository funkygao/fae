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

// specs: inbound params use json
func (this *FunServantImpl) unmarshalIn(d []byte) (v bson.M, err error) {
	err = json.Unmarshal(d, &v)
	if err != nil {
		log.Error("unmarshalIn error: %s -> %s", d, err)
	}

	return
}

// specs: outbound data use bson
func (this *FunServantImpl) marshalOut(d bson.M) []byte {
	val, err := bson.Marshal(d)
	if err != nil {
		// should never happen
		log.Critical("marshalOut error: %+v -> %v", d, err)
	}
	return val
}

func (this *FunServantImpl) mgFieldsIsNil(fields []byte) bool {
	return len(fields) <= 5
}
