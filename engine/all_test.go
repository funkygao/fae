package engine

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestEngineConfig(t *testing.T) {
	e := NewEngine()
	e.LoadConfigFile("../etc/fproxyd.cf")
	assert.Equal(t, ":9001", e.conf.listenAddr)
	assert.NotEqual(t, 0, len(e.conf.memcaches))
	assert.NotEqual(t, 0, len(e.conf.mongos))
}
