package peer

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestPeerMessage(t *testing.T) {
	var msg = peerMessage{}
	msg["cmd"] = "ok"
	data, _ := msg.marshal()
	assert.Equal(t, `{"cmd":"ok"}`, string(data))
}
