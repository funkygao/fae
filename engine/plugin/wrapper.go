package plugin

import (
	conf "github.com/funkygao/jsconf"
)

// A helper object to support delayed plugin creation.
type pluginWrapper struct {
	name          string
	configCreator func() *conf.Conf
	pluginCreator func() Plugin
}

func (this *pluginWrapper) Create() (plugin Plugin) {
	plugin = this.pluginCreator()
	plugin.Init(this.configCreator())
	return
}
