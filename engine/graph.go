package engine

import (
	"html/template"
	"io"
	"sync"
	"time"
)

type graphPoints [2]int

type Graph struct {
	Title                                                   string
	Qps, ActiveSessions, Latencies, Errors, Calls, Sessions []graphPoints
	Tpl                                                     *template.Template
	mu                                                      sync.Mutex
	rpcServer                                               *TFunServer
}

func NewGraph(title, tpl string, rpcServer *TFunServer) Graph {
	return Graph{
		Title: title,
		Tpl:   template.Must(template.New("vis").Parse(tpl)),

		Qps:            []graphPoints{},
		ActiveSessions: []graphPoints{},
		Latencies:      []graphPoints{},
		Errors:         []graphPoints{},
		Calls:          []graphPoints{},
		Sessions:       []graphPoints{},

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
		g.Calls = []graphPoints{}
		g.Sessions = []graphPoints{}

	}

	ts := int(time.Now().UnixNano() / 1e6)
	g.Qps = append(g.Qps, graphPoints{ts,
		int(g.rpcServer.stats.CallPerSecond.Rate1())})
	g.ActiveSessions = append(g.ActiveSessions, graphPoints{ts,
		int(g.rpcServer.activeSessionN)})
	g.Latencies = append(g.Latencies, graphPoints{ts,
		int(g.rpcServer.stats.CallLatencies.Mean())})
	g.Errors = append(g.Errors, graphPoints{ts,
		int(g.rpcServer.cumCallErrs)})
	g.Calls = append(g.Calls, graphPoints{ts,
		int(g.rpcServer.cumCalls)})
	g.Sessions = append(g.Sessions, graphPoints{ts,
		int(g.rpcServer.cumSessions)})

	g.Tpl.Execute(w, g)
}
