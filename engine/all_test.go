package engine

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestEngineConfig(t *testing.T) {
	e := NewEngine()
	e.LoadConfigFile("../etc/faed.cf")
	assert.Equal(t, ":9001", e.conf.rpc.listenAddr)
	assert.NotEqual(t, 0, len(e.conf.memcaches))
	assert.NotEqual(t, 0, len(e.conf.mongos))
}
