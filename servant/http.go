package servant

import (
	"fmt"
	"github.com/funkygao/golib/server"
	"github.com/gorilla/mux"
	"net/http"
)

func (this *FunServantImpl) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		cmd    = vars["cmd"]
		output = make(map[string]interface{})
	)

	switch cmd {
	case "stat", "stats":
		if this.mg != nil {
			output["mongo"] = this.mg.FreeConnMap()
		}
		if this.mc != nil {
			output["memcache"] = this.mc.FreeConnMap()
		}
		if this.lc != nil {
			output["lcache"] = this.lc.Len()
		}
		if this.proxy != nil {
			output["proxy"] = this.proxy.StatsMap()
		}

		calls := make(map[string]interface{})
		for _, key := range svtStats.calls.Keys() {
			calls[key] = fmt.Sprintf("%.2f%%", svtStats.calls.Percent(key))
		}
		output["rpc.call"] = calls
		output["runtime"] = this.Runtime()

	case "conf":
		output["conf"] = *this.conf

	case "guide", "help", "h":
		output["uris"] = []string{
			"/svt/stat",
			"/svt/conf",
		}

	default:
		return nil, server.ErrHttp404
	}

	return output, nil
}
