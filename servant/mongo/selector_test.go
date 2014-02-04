package mongo

import (
	"github.com/funkygao/assert"
	"github.com/funkygao/fae/config"
	conf "github.com/funkygao/jsconf"
	"testing"
)

func setupConfig() *config.ConfigMongodb {
	cf, _ := conf.Load("../../etc/faed.cf")
	section, _ := cf.Section("servants")
	config.LoadServants(section)
	return config.Servants.Mongodb
}

func TestPickServer(t *testing.T) {
	cf := setupConfig()
	picker := NewStandardServerSelector(1000)
	picker.SetServers(cf.Servers)

	addr, err := picker.PickServer("db", 23)
	assert.Equal(t, "mongodb://127.0.0.1:27017/", addr)
	assert.Equal(t, nil, err)

	addr, err = picker.PickServer("invalid", 23)
	assert.Equal(t, "", addr)
	assert.Equal(t, ErrServerNotFound, err)

	addr, err = picker.PickServer("db", 2300) // too big for 1000
	assert.Equal(t, "", addr)
	assert.Equal(t, ErrServerNotFound, err)

	addr, err = picker.PickServer("default", 1<<30) // too big for 1000
	assert.Equal(t, "mongodb://127.0.0.1:27017/", addr)
	assert.Equal(t, nil, err)

	addr, err = picker.PickServer("log", 1<<30) // too big for 1000
	assert.Equal(t, "mongodb://127.0.0.1:27017/", addr)
	assert.Equal(t, nil, err)
}
