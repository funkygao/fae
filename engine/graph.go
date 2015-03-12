package engine

import (
	"html/template"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type graphPoints [2]int

type Graph struct {
	Title                                         string
	Qps, ActiveSessions, Latencies, Errors, Slows []graphPoints
	Calls, Sessions                               int64
	Tpl                                           *template.Template
	mu                                            sync.Mutex
	rpcServer                                     *TFunServer
}

func NewGraph(title, tpl string, rpcServer *TFunServer) Graph {
	return Graph{
		Title: title,
		Tpl:   template.Must(template.New("vis").Parse(tpl)),

		Qps:            []graphPoints{},
		ActiveSessions: []graphPoints{},
		Latencies:      []graphPoints{},
		Errors:         []graphPoints{},
		Slows:          []graphPoints{},

		rpcServer: rpcServer,
	}
}

func (g *Graph) write(w io.Writer) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.Qps) > (1 << 20) {
		// dashboard should never use up fae's memory
		g.Qps = []graphPoints{}
		g.Latencies = []graphPoints{}
		g.ActiveSessions = []graphPoints{}
		g.Errors = []graphPoints{}
		g.Slows = []graphPoints{}

	}

	ts := int(time.Now().UnixNano() / 1e6)
	g.Qps = append(g.Qps, graphPoints{ts,
		int(g.rpcServer.stats.CallPerSecond.Rate1())})
	g.ActiveSessions = append(g.ActiveSessions, graphPoints{ts,
		int(g.rpcServer.activeSessionN)})
	g.Latencies = append(g.Latencies, graphPoints{ts,
		int(g.rpcServer.stats.CallLatencies.Mean())})
	errs := atomic.LoadInt64(&g.rpcServer.cumCallErrs)
	g.Errors = append(g.Errors, graphPoints{ts,
		int(errs)})
	slows := atomic.LoadInt64(&g.rpcServer.cumCallSlow)
	g.Slows = append(g.Slows, graphPoints{ts,
		int(slows)})

	g.Calls = atomic.LoadInt64(&g.rpcServer.cumCalls)
	g.Sessions = atomic.LoadInt64(&g.rpcServer.cumSessions)

	g.Tpl.Execute(w, g)
}
