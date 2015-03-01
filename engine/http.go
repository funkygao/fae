package engine

import (
	"fmt"
	"github.com/funkygao/fae/config"
	"github.com/funkygao/golib/server"
	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
	"net/http"
	"runtime"
	"syscall"
	"time"
)

func (this *Engine) stopHttpServ() {
	server.StopHttpServ()
}

func (this *Engine) launchHttpServ() {
	if config.Engine.HttpListenAddr == "" {
		return
	}

	server.LaunchHttpServ(config.Engine.HttpListenAddr, config.Engine.PprofListenAddr)
	server.RegisterHttpApi("/h", func(w http.ResponseWriter, req *http.Request,
		params map[string]interface{}) (interface{}, error) {
		return this.handleHttpHelpQuery(w, req, params)
	}).Methods("GET")
	server.RegisterHttpApi("/help", func(w http.ResponseWriter, req *http.Request,
		params map[string]interface{}) (interface{}, error) {
		return this.handleHttpHelpQuery(w, req, params)
	}).Methods("GET")

	server.RegisterHttpApi("/engine/{cmd}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *Engine) handleHttpHelpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	output := make(map[string]interface{})
	if config.Engine.PprofListenAddr != "" {
		output["pprof"] = "http://" + config.Engine.PprofListenAddr + "/debug/pprof/"
		output["vars"] = "http://" + config.Engine.PprofListenAddr + "/debug/vars"
	}

	output["uris"] = []string{"/engine/help", "/svt/help"}
	return output, nil
}

func (this *Engine) handleHttpQuery(w http.ResponseWriter, req *http.Request,
	params map[string]interface{}) (interface{}, error) {
	var (
		vars   = mux.Vars(req)
		cmd    = vars["cmd"]
		output = make(map[string]interface{})
	)

	switch cmd {
	case "ping":
		output["status"] = "ok"

	case "debug":
		stack := make([]byte, 1<<20)
		stackSize := runtime.Stack(stack, true)
		output["result"] = "go to global logger to see result"
		log.Info(string(stack[:stackSize]))

	case "stop":
		this.rpcServer.Stop()
		output["status"] = "stopped"

	case "stat", "stats":
		rusage := syscall.Rusage{}
		syscall.Getrusage(0, &rusage)
		output["rusage"] = rusage
		output["started"] = this.StartedAt
		output["elapsed"] = time.Since(this.StartedAt).String()
		output["pid"] = this.pid
		output["hostname"] = this.hostname
		output["ver"] = server.VERSION
		output["build_id"] = server.BuildID
		output["active_sessions"] = 1 // TODO

	case "qps":
		output["qps"] = fmt.Sprintf("1m:%.0f, 5m:%.0f 15m:%.0f avg:%.0f",
			this.stats.CallPerSecond.Rate1(),
			this.stats.CallPerSecond.Rate5(),
			this.stats.CallPerSecond.Rate15(),
			this.stats.CallPerSecond.RateMean())

	case "runtime":
		output["runtime"] = this.stats.Runtime()

	case "mem":
		output["mem"] = *this.stats.memStats

	case "conf":
		output["engine"] = *config.Engine
		output["rpc"] = *config.Engine.Rpc
		output["servants"] = *config.Engine.Servants

	case "guide", "help", "h":
		output["uris"] = []string{
			"/engine/ping",
			"/engine/debug",
			"/engine/stop",
			"/engine/stat",
			"/engine/qps",
			"/engine/runtime",
			"/engine/mem",
			"/engine/conf",
		}

	default:
		return nil, server.ErrHttp404
	}

	return output, nil
}
