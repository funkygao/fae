package plugins

import (
	"github.com/funkygao/fae/engine"
	conf "github.com/funkygao/jsconf"
)

// Just a demo of how to use plugin mechanism.
type PluginFoo struct {
}

func (this *PluginFoo) Init(cf *conf.Conf) {

}

func init() {
	engine.RegisterPlugin("foo", func() engine.Plugin {
		return new(PluginFoo)
	})

}
