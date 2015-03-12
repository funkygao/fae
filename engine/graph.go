package engine

import (
	"html/template"
	"io"
	"sync"
	"time"
)

type graphPoints [2]int

type Graph struct {
	Title                                                             string
	HeapUse, ScvgInuse, ScvgIdle, ScvgSys, ScvgReleased, ScvgConsumed []graphPoints
	Tpl                                                               *template.Template
	mu                                                                sync.RWMutex
	rpcServer                                                         *TFunServer
}

func NewGraph(title, tpl string, rpcServer *TFunServer) Graph {
	return Graph{
		Title:        title,
		HeapUse:      []graphPoints{},
		ScvgInuse:    []graphPoints{},
		ScvgIdle:     []graphPoints{},
		ScvgSys:      []graphPoints{},
		ScvgReleased: []graphPoints{},
		ScvgConsumed: []graphPoints{},
		rpcServer:    rpcServer,
		Tpl:          template.Must(template.New("vis").Parse(tpl)),
	}
}

func (g *Graph) write(w io.Writer) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	ts := int(time.Now().UnixNano() / 1e6)
	g.HeapUse = append(g.HeapUse, graphPoints{ts,
		int(g.rpcServer.stats.CallPerSecond.Rate1())})

	g.Tpl.Execute(w, g)
}
