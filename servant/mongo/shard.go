package mongo

import (
	"fmt"
	"github.com/funkygao/fxi/config"
)

func lookupDbName(shardKey string, shardId int) string {
	n := (shardId / config.Servants.Mongodb.ShardBaseNum) + 1
	return fmt.Sprintf("db%s", n)
}
