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

		if this.game != nil {
			h := this.game.PhpLatencyHistogram().Snapshot()
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99})
			latency := make(map[string]interface{})
			latency["count"] = h.Count()
			latency["min"] = h.Min()
			latency["max"] = h.Max()
			latency["mean"] = h.Mean()
			latency["stdev"] = h.StdDev()
			latency["median"] = ps[0]
			latency["75%"] = ps[1]
			latency["95%"] = ps[2]
			latency["99%"] = ps[3]
			output["php.latency"] = latency

			h = this.game.PhpPayloadHistogram().Snapshot()
			ps = h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99})
			payload := make(map[string]interface{})
			payload["count"] = h.Count()
			payload["min"] = h.Min()
			payload["max"] = h.Max()
			payload["mean"] = h.Mean()
			payload["stdev"] = h.StdDev()
			payload["median"] = ps[0]
			payload["75%"] = ps[1]
			payload["95%"] = ps[2]
			payload["99%"] = ps[3]
			output["php.payload"] = payload
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
