package game

import (
	"github.com/funkygao/fae/config"
	"testing"
)

func BenchmarkRegister(b *testing.B) {
	b.ReportAllocs()
	cf := &config.ConfigGame{
		RedisServerAddr: "127.0.0.1:6379",
		ShardSplit:      config.ConfigGameShardSplit{User: 1 << 30},
	}
	r := newRegister(cf)
	for i := 0; i < b.N; i++ {
		r.Register(REG_USER)
	}
}
