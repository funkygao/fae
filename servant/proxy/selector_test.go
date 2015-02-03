package proxy

import (
	"testing"
)

func TestRandServant(t *testing.T) {
	s := newStandardPeerSelector()
	s.SetPeersAddr([]string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}...)
	for i := 0; i < 5; i++ {
		t.Logf("%s", s.RandPeer())
	}
}
