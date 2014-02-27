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
		client := prepareFixture(i)
		t.Logf("members %d", i)
		deleteDeadSessions(t, client)
	}

}

func deleteDeadSessions(t *testing.T, client *mockClient) {
	for _, sessions := range client.freeconn {
		for idx, _ := range sessions {
			t.Logf("idx %d", idx)
			if rand.Intn(2) > -1 {
				t.Logf("before delete: %+v", client.freeconn["localhost"])
				t.Logf("idx %d", idx)
				deleteSession(client, sessions, idx)
				t.Logf("after  delete: %+v", client.freeconn["localhost"])
			}
		}
	}

}

func deleteSession(c *mockClient, sessions []*mockSession, idx int) {
	//sessions = append(sessions[:idx], sessions[idx+1:]...) this doesn't work
	c.freeconn["localhost"] = append(c.freeconn["localhost"][:idx],
		c.freeconn["localhost"][idx+1:]...)

}
