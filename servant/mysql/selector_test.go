package mysql

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	"testing"
)

func TestSelectorStandardEndsWithDigit(t *testing.T) {
	s := newStandardServerSelector(new(config.ConfigMysql))
	assert.Equal(t, true, s.endsWithDigit("AllianceShard8"))
	assert.Equal(t, false, s.endsWithDigit("ShardLookup"))
}

func TestSelectorStandardPoolServers(t *testing.T) {
	s := server.NewServer("test")
	s.LoadConfig("../../etc/faed.cf.sample")
	section, _ := s.Conf.Section("servants.mysql")
	cf := &config.ConfigMysql{}
	cf.LoadConfig(section)

	sel := newStandardServerSelector(cf)

	assert.Equal(t, 1, len(sel.PoolServers("UserShard")))
}
