package mongo

import (
	"fmt"
	"net"
)

type ServerSelector interface {
	SetServers(servers ...string) error
	PickServer(key string) (net.Addr, error)
}

type ShardServerSelector struct {
	ShardBaseNum int
}

func (this *ShardServerSelector) PickServer(shardKey string, shardId int) string {
	n := (shardId / this.ShardBaseNum) + 1
	return fmt.Sprintf("db%s", n)
}

func (this *ShardServerSelector) SetServers() error {
	return nil
}
