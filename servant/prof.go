package servant

import (
	"github.com/funkygao/golib/sampling"
	"time"
)

type profilerInfo struct {
	do bool
	t1 time.Time
}

func (this *FunServantImpl) profilerInfo() profilerInfo {
	info := profilerInfo{do: false}
	info.do = sampling.SampleRateSatisfied(this.conf.ProfilerRate)
	if info.do {
		info.t1 = time.Now()
	}

	return info
}

func (this *FunServantImpl) truncValue(val []byte) []byte {
	if len(val) < this.conf.ProfilerMaxAnswerSize {
		return val
	}

	return val[:this.conf.ProfilerMaxAnswerSize]
}
