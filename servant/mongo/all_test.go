package mongo

import (
	"math/rand"
	"testing"
)

type mockSession struct {
	id int
}

type mockClient struct {
	freeconn map[string][]*mockSession
}

var client *mockClient

func prepareFixture(n int) *mockClient {
	addr := "localhost"
	c := &mockClient{freeconn: make(map[string][]*mockSession)}
	for i := 0; i < n; i++ {
		c.freeconn[addr] = append(c.freeconn[addr], &mockSession{id: i})
	}
	return c
}

func TestDeleteDeadSessions(t *testing.T) {
	for i := 0; i < 10; i++ {
		client = prepareFixture(i)
		t.Logf("members %d", i)
		testDeleteDeadSessions(t)
	}

}

func testDeleteDeadSessions(t *testing.T) {
	for _, sessions := range client.freeconn {
		for _, sess := range sessions {
			if rand.Intn(2) == 1 {
				t.Logf("before delete: %+v", client.freeconn["localhost"])
				deleteSession(sess)
				t.Logf("after  delete: %+v", client.freeconn["localhost"])
			}
		}
	}

}

func deleteSession(session *mockSession) {
	//sessions = append(sessions[:idx], sessions[idx+1:]...) this doesn't work
	for idx, sess := range client.freeconn["localhost"] {
		if sess == session {
			client.freeconn["localhost"] = append(client.freeconn["localhost"][:idx],
				client.freeconn["localhost"][idx+1:]...)
		}
	}

}
