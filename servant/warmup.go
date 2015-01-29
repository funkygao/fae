package servant

import (
	log "github.com/funkygao/log4go"
)

func (this *FunServantImpl) warmUp() {
	log.Debug("warming up...")

	if this.mg != nil {
		go this.mg.Warmup()
	}

	if this.mc != nil {
		go this.mc.Warmup()
	}

	if this.my != nil {
		this.my.Warmup()
	}

	if this.proxy != nil {
		this.proxy.Warmup()
	}

	log.Debug("warmup done")
}
