package engine

import (
	"github.com/funkygao/etclib"
	"github.com/funkygao/fae/config"
	"html/template"
	"io"
	"math"
	"net"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type graphPoints [2]int

type uint64Slice []uint64

func (s uint64Slice) Len() int {
	return len(s)
}

func (s uint64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s uint64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

// Dashboard data of engine.
type graph struct {
	Title                                         string
	Qps, ActiveSessions, Latencies, Errors, Slows []graphPoints
	NumGC, HeapSys, HeapAlloc, HeapReleased       []graphPoints
	StackInUse                                    []graphPoints
	HeapObjects                                   []graphPoints
	GcPause100, GcPause99, GcPause95, GcPause75   []graphPoints
	Calls, Sessions                               int64
	Peers                                         []string
	Tpl                                           *template.Template
	mu                                            sync.Mutex
	rpcServer                                     *TFunServer
	port                                          string
}

func newGraph(title, tpl string, rpcServer *TFunServer) graph {
	_, port, _ := net.SplitHostPort(config.Engine.DashboardListenAddr)
	return graph{
		Title:          title,
		Tpl:            template.Must(template.New("vis").Parse(tpl)),
		port:           port,
		Qps:            []graphPoints{},
		ActiveSessions: []graphPoints{},
		Latencies:      []graphPoints{},
		Errors:         []graphPoints{},
		Slows:          []graphPoints{},
		NumGC:          []graphPoints{},
		HeapSys:        []graphPoints{},
		HeapAlloc:      []graphPoints{},
		HeapReleased:   []graphPoints{},
		StackInUse:     []graphPoints{},
		HeapObjects:    []graphPoints{},
		GcPause100:     []graphPoints{},
		GcPause99:      []graphPoints{},
		GcPause95:      []graphPoints{},
		GcPause75:      []graphPoints{},
		rpcServer:      rpcServer,
	}
}

func (g *graph) write(w io.Writer) {
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

	if len(g.Qps) > (60 * 60 * 12 / 10) { // 12 hours
		// dashboard should never use up fae's memory
		g.Qps = []graphPoints{}
		g.Latencies = []graphPoints{}
		g.ActiveSessions = []graphPoints{}
		g.Errors = []graphPoints{}
		g.Slows = []graphPoints{}
		g.HeapAlloc = []graphPoints{}
		g.HeapReleased = []graphPoints{}
		g.HeapSys = []graphPoints{}
		g.NumGC = []graphPoints{}
		g.StackInUse = []graphPoints{}
		g.HeapObjects = []graphPoints{}
		g.GcPause100 = []graphPoints{}
		g.GcPause99 = []graphPoints{}
		g.GcPause95 = []graphPoints{}
		g.GcPause75 = []graphPoints{}
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

	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	g.NumGC = append(g.NumGC, graphPoints{ts,
		int(memStats.NumGC)})
	g.HeapSys = append(g.HeapSys, graphPoints{ts,
		int(memStats.HeapSys) / (1 << 20)})
	g.HeapReleased = append(g.HeapReleased, graphPoints{ts,
		int(memStats.HeapReleased) / (1 << 20)})
	g.HeapAlloc = append(g.HeapAlloc, graphPoints{ts,
		int(memStats.HeapAlloc) / (1 << 20)})
	g.StackInUse = append(g.StackInUse, graphPoints{ts,
		int(memStats.StackInuse) / (1 << 20)})
	g.HeapObjects = append(g.HeapObjects, graphPoints{ts,
		int(memStats.HeapObjects)})

	// sort the GC pause array
	gcPausesMs := make(uint64Slice, 0, len(memStats.PauseNs))
	for _, pauseNs := range memStats.PauseNs {
		if pauseNs == 0 {
			continue
		}

		pauseMs := pauseNs / uint64(time.Millisecond)
		if pauseMs == 0 {
			continue
		}

		gcPausesMs = append(gcPausesMs, pauseMs)
	}
	sort.Sort(gcPausesMs)

	g.GcPause100 = append(g.GcPause100, graphPoints{ts,
		int(percentile(100.0, gcPausesMs))})
	g.GcPause99 = append(g.GcPause99, graphPoints{ts,
		int(percentile(99.0, gcPausesMs))})
	g.GcPause95 = append(g.GcPause95, graphPoints{ts,
		int(percentile(95.0, gcPausesMs))})
	g.GcPause75 = append(g.GcPause75, graphPoints{ts,
		int(percentile(75.0, gcPausesMs))})

	g.Tpl.Execute(w, g)
}

func percentile(perc float64, sortedArr []uint64) uint64 {
	length := len(sortedArr)
	if length == 0 {
		return 0
	}
	indexOfPerc := int(math.Floor(((perc / 100.0) * float64(length)) + 0.5))
	if indexOfPerc >= length {
		indexOfPerc = length - 1
	}
	return sortedArr[indexOfPerc]
}
