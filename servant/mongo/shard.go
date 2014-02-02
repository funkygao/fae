package mongo

import (
	"fmt"
	"github.com/funkygao/fae/config"
)

func lookupDbName(shardKey string, shardId int) string {
	n := (shardId / config.Servants.Mongodb.ShardBaseNum) + 1
	return fmt.Sprintf("db%s", n)
}
