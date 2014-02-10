package engine

import (
	"encoding/json"
	"errors"
	log "github.com/funkygao/log4go"
	"github.com/gorilla/mux"
	"io"
	"net"
	"net/http"
	"runtime"
	"time"
)

func (this *Engine) setupHttpServ() {
	if this.conf.httpListenAddr == "" {
		return
	}

	this.httpRouter = mux.NewRouter()
	this.httpServer = &http.Server{Addr: this.conf.httpListenAddr,
		Handler: this.httpRouter}

	this.RegisterHttpApi("/admin/{cmd}",
		func(w http.ResponseWriter, req *http.Request,
			params map[string]interface{}) (interface{}, error) {
			return this.handleHttpQuery(w, req, params)
		}).Methods("GET")
}

func (this *Engine) launchHttpServ() {
	var err error
	this.httpListener, err = net.Listen("tcp", this.httpServer.Addr)
	if err != nil {
		panic(err)
	}

	log.Info("HTTP server ready at %s", this.conf.httpListenAddr)

	go this.httpServer.Serve(this.httpListener)
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

	case "reload":
		this.LoadConfigFile()
		output["status"] = "ok"

	case "stop":
		this.rpcServer.Stop()
		output["status"] = "stopped"

	case "stat":
		output["started"] = this.StartedAt
		output["elapsed"] = time.Since(this.StartedAt).String()
		output["pid"] = this.pid
		output["hostname"] = this.hostname
		output["stats"] = this.stats
		output["peers"] = this.peer.Neighbors()

	case "runtime":
		output["runtime"] = this.stats.Runtime()

	case "mem":
		output["mem"] = *this.stats.memStats

	case "uris":
		output["all"] = this.httpPaths

	default:
		return nil, errors.New("Not Found")
	}

	return output, nil
}

func (this *Engine) RegisterHttpApi(path string,
	handlerFunc func(http.ResponseWriter,
		*http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunc := func(w http.ResponseWriter, req *http.Request) {
		var (
			ret interface{}
			t1  = time.Now()
		)

		params, err := this.decodeHttpParams(w, req)
		if err == nil {
			ret, err = handlerFunc(w, req, params)
		}

		if err != nil {
			ret = map[string]interface{}{"error": err.Error()}
		}

		w.Header().Set("Content-Type", "application/json")
		var status int
		if err == nil {
			status = http.StatusOK
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

		// debug request body content
		log.Trace("req body: %+v", params)
		// access log
		log.Debug("%s \"%s %s %s\" %d %s",
			req.RemoteAddr,
			req.Method,
			req.RequestURI,
			req.Proto,
			status,
			time.Since(t1))
		if status != http.StatusOK {
			log.Error("ERROR %v", err)
		}

		if ret != nil {
			// pretty write json result
			pretty, _ := json.MarshalIndent(ret, "", "    ")
			w.Write(pretty)
			w.Write([]byte("\n"))
		}
	}

	// path can't be duplicated
	isDup := false
	for _, p := range this.httpPaths {
		if p == path {
			log.Error("REST[%s] already registered", path)
			isDup = true
			break
		}
	}

	if !isDup {
		this.httpPaths = append(this.httpPaths, path)
	}

	return this.httpRouter.HandleFunc(path, wrappedFunc)
}

func (this *Engine) decodeHttpParams(w http.ResponseWriter,
	req *http.Request) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return params, nil
}

func (this *Engine) stopHttpServ() {
	if this.httpListener != nil {
		this.httpListener.Close()

		log.Info("HTTP server stopped")
	}
}
