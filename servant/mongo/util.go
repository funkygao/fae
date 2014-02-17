package mongo

import (
	log "github.com/funkygao/log4go"
	"labix.org/v2/mgo/bson"
)

// specs: inbound params use json
func UnmarshalIn(d []byte) (v bson.M, err error) {
	err = bson.Unmarshal(d, &v)
	if err != nil {
		log.Error("mongo.unmarshalIn error: %s -> %s", d, err)
	}

	return
}

// specs: outbound data use bson
func MarshalOut(d bson.M) []byte {
	val, err := bson.Marshal(d)
	if err != nil {
		// should never happen
		log.Critical("mongo.marshalOut error: %+v -> %v", d, err)
	}

	return val
}

func FieldsIsNil(fields []byte) bool {
	return len(fields) <= 5
}
