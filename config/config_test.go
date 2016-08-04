package config

import (
	"encoding/json"
	"testing"

	"github.com/funkygao/golib/server"
)

func TestJsonizeConfigMysql(t *testing.T) {
	s := server.NewServer("test")
	s.LoadConfig("../etc/pubsub.cf")
	section, _ := s.Conf.Section("servants.mysql")
	cf := &ConfigMysql{}
	cf.LoadConfig(section)

	j, err := json.Marshal(cf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", string(j))
}
