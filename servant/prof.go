package servant

import (
	"github.com/funkygao/fae/servant/gen-go/fun/rpc"
	"github.com/funkygao/golib/sampling"
	"time"
)

type profilerInfo struct {
	do bool
	t1 time.Time
}

func (this *FunServantImpl) profilerInfo(ctx *rpc.Context) profilerInfo {
	var sampleRate int16 = 1000
	if ctx.IsSetProfRate() {
		sampleRate = *ctx.ProfRate
	}

	info := profilerInfo{do: false}
	info.do = sampling.SampleRateSatisfied(int(sampleRate))
	if info.do {
		info.t1 = time.Now()
	}

	return info
}
