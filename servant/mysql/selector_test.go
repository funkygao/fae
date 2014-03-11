package mysql

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	"testing"
)

func TestSelectorStandardEndsWithDigit(t *testing.T) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	assert.Equal(t, true, s.endsWithDigit("AllianceShard8"))
	assert.Equal(t, false, s.endsWithDigit("ShardLookup"))
}

func BenchmarkEndsWithDigit(b *testing.B) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	for i := 0; i < b.N; i++ {
		s.endsWithDigit("UserShard1")
	}
}
