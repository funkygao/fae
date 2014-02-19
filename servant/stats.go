package servant

import (
	"github.com/funkygao/metrics"
)

type servantStats struct {
	Ping   metrics.Counter
	IdNext metrics.Counter
	Dlog   metrics.Counter

	LcSet metrics.Counter
	LcGet metrics.Counter
	LcDel metrics.Counter

	McSet metrics.Counter
	McGet metrics.Counter
	McAdd metrics.Counter
	McDel metrics.Counter
	McInc metrics.Counter

	MgFindOne       metrics.Counter
	MgFindAll       metrics.Counter
	MgCount         metrics.Counter
	MgUpdate        metrics.Counter
	MgUpsert        metrics.Counter
	MgInsert        metrics.Counter
	MgInserts       metrics.Counter
	MgDel           metrics.Counter
	MgFindAndModify metrics.Counter
}

func (this *servantStats) registerMetrics() {
	this.Ping = metrics.NewCounter()
	metrics.Register("s.ping", this.Ping)
	this.IdNext = metrics.NewCounter()
	metrics.Register("s.id.next", this.IdNext)
	this.Dlog = metrics.NewCounter()
	metrics.Register("s.dlog", this.Dlog)

	this.LcSet = metrics.NewCounter()
	metrics.Register("s.lc.set", this.LcSet)
	this.LcGet = metrics.NewCounter()
	metrics.Register("s.lc.get", this.LcGet)
	this.LcDel = metrics.NewCounter()
	metrics.Register("s.lc.del", this.LcDel)

	this.McAdd = metrics.NewCounter()
	metrics.Register("s.mc.add", this.McAdd)
	this.McDel = metrics.NewCounter()
	metrics.Register("s.mc.del", this.McDel)
	this.McGet = metrics.NewCounter()
	metrics.Register("s.mc.get", this.McGet)
	this.McSet = metrics.NewCounter()
	metrics.Register("s.mc.set", this.McSet)
	this.McInc = metrics.NewCounter()
	metrics.Register("s.mc.inc", this.McInc)

	this.MgFindAll = metrics.NewCounter()
	metrics.Register("s.mg.findAll", this.MgFindAll)
	this.MgFindOne = metrics.NewCounter()
	metrics.Register("s.mg.findOne", this.MgFindOne)
	this.MgFindAndModify = metrics.NewCounter()
	metrics.Register("s.mg.findAndModify", this.MgFindAndModify)
	this.MgCount = metrics.NewCounter()
	metrics.Register("s.mg.count", this.MgCount)
	this.MgDel = metrics.NewCounter()
	metrics.Register("s.mg.del", this.MgDel)
	this.MgUpdate = metrics.NewCounter()
	metrics.Register("s.mg.update", this.MgUpdate)
	this.MgUpsert = metrics.NewCounter()
	metrics.Register("s.mg.upsert", this.MgUpsert)
	this.MgInsert = metrics.NewCounter()
	metrics.Register("s.mg.insert", this.MgInsert)
	this.MgInserts = metrics.NewCounter()
	metrics.Register("s.mg.inserts", this.MgInserts)
}
