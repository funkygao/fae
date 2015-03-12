package engine

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"html/template"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type graphPoints [2]int

type Graph struct {
	Title                                         string
	Qps, ActiveSessions, Latencies, Errors, Slows []graphPoints
	Calls, Sessions                               int64
	Peers                                         []string
	Tpl                                           *template.Template
	mu                                            sync.Mutex
	rpcServer                                     *TFunServer
	port                                          string
}

func NewGraph(title, tpl string, rpcServer *TFunServer) Graph {
	_, port, _ := net.SplitHostPort(config.Engine.DashboardListenAddr)
	return Graph{
		Title:          title,
		Tpl:            template.Must(template.New("vis").Parse(tpl)),
		port:           port,
		Qps:            []graphPoints{},
		ActiveSessions: []graphPoints{},
		Latencies:      []graphPoints{},
		Errors:         []graphPoints{},
		Slows:          []graphPoints{},
		rpcServer:      rpcServer,
	}
}

func (g *Graph) write(w io.Writer) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// get all fae dashboard url in cluster
	g.Peers = make([]string, 0, 10)
	if peers, err := etclib.ServiceEndpoints(etclib.SERVICE_FAE); err == nil && len(peers) > 0 {
		for _, peer := range peers {
			host, _, _ := net.SplitHostPort(peer)
			g.Peers = append(g.Peers, host+":"+g.port)
		}
	}

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
