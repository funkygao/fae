package engine

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestEngineConfig(t *testing.T) {
	e := NewEngine("../etc/faed.cf")
	e.LoadConfigFile()
	assert.Equal(t, ":9001", e.conf.rpc.listenAddr)
}

func TestPeerMessage(t *testing.T) {
	var msg = peerMessage{}
	msg["cmd"] = "ok"
	data, _ := msg.marshal()
	assert.Equal(t, `{"cmd":"ok"}`, string(data))
}
