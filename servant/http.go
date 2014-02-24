package servant

import (
	rest "github.com/funkygao/fae/http"
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
	case "stat":
		if this.mg != nil {
			output["mongo"] = this.mg.FreeConn()
		}
		if this.mc != nil {
			output["memcache"] = this.mc.FreeConn()
		}
		if this.lc != nil {
			output["lcache"] = this.lc.Len()
		}
		if this.proxy != nil {
			output["proxy"] = this.proxy.StatsJSON()
		}

	case "conf":
		output["conf"] = *this.conf

	default:
		return nil, rest.ErrHttp404
	}

	return output, nil
}
