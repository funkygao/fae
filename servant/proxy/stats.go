package proxy

import (
	"bytes"
	"fmt"
)

func (this *Proxy) StatsJSON() string {
	s := new(bytes.Buffer)
	for addr, pool := range this.pools {
		s.WriteString(fmt.Sprintf(`{"%s":%v}`, addr, pool.pool.StatsJSON()))
	}
	return s.String()
}
