package mongo

import (
	"github.com/funkygao/fae/config"
	"labix.org/v2/mgo"
)

type MongodbPool struct {
	sessions map[string]*mgo.Session // key is shard name
}

func (this *MongodbPool) Init(cf *config.ConfigMongodb) {

}
