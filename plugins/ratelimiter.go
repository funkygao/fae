package plugins

import (
	"github.com/funkygao/fae/engine/plugin"
	conf "github.com/funkygao/jsconf"
	log "github.com/funkygao/log4go"
)

type RateLimiter struct {
	limit int
}

func (this *RateLimiter) Init(cf *conf.Conf) {
	this.limit = cf.Int("limit", 2000)

	log.Debug("%T init: %+v", *this, *this)
}

func (this *RateLimiter) Run(r plugin.PluginRunner) error {
	log.Debug("%T started", *this)

	var (
		inChan = r.InChan()
		pack   *plugin.PipelinePack
	)

	for {
		pack = <-inChan

		r.Inject(pack)
	}
	return nil
}

func init() {
	plugin.RegisterPlugin("RateLimiter", func() plugin.Plugin {
		return new(RateLimiter)
	})

}
