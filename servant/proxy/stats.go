package proxy

import (
	"bytes"
	"fmt"
)

func (this *Proxy) StatsJSON() string {
	s := new(bytes.Buffer)
	s.WriteString("[\n")
	for addr, pool := range this.pools {
		s.WriteString(fmt.Sprintf(`%s{"%s":%v}%s`, "\t",
			addr, pool.pool.StatsJSON(), "\n"))
	}
	s.WriteString("]")
	return s.String()
}
