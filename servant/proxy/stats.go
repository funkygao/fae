package proxy

import (
	"encoding/json"
)

func (this *Proxy) StatsJSON() string {
	m := make(map[string]string)
	for addr, pool := range this.pools {
		m[addr] = pool.pool.StatsJSON()
	}

	pretty, _ := json.MarshalIndent(m, "", "    ")
	return string(pretty)
}

func (this *Proxy) StatsMap() map[string]string {
	m := make(map[string]string)
	for addr, pool := range this.pools {
		m[addr] = pool.pool.StatsJSON()
	}

	return m
}
