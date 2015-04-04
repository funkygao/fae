package engine

import (
	"github.com/funkygao/fae/config"
	"net/http"
)

func (this *Engine) launchDashboard() {
	if config.Engine.DashboardListenAddr == "" {
		return
	}

	this.graph = newGraph("RPC Dashboard", dashboard_tpl, this.rpcServer.(*TFunServer))

	http.HandleFunc("/", this.dashboard)
	go http.ListenAndServe(config.Engine.DashboardListenAddr, nil)
}

func (this *Engine) dashboard(w http.ResponseWriter, r *http.Request) {
	this.graph.write(w)
}
