package engine

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
)

var (
	registeredPlugins = make(map[string]func() Plugin) // name:factory
)

type Plugin interface {
	Init(section *conf.Conf)
}

func RegisterPlugin(name string, factory func() Plugin) {
	if _, present := registeredPlugins[name]; present {
		panic(fmt.Sprintf("plugin[%s] cannot register twice", name))
	}

	registeredPlugins[name] = factory
}

// All plugins have these attributes.
type pluginCommon struct {
	name    string `json:"name"`
	class   string `json:"class"`
	enabled bool   `json:"enabled"`
}

func (this *pluginCommon) load(section *conf.Conf) {
	this.name = section.String("name", "")
	if this.name == "" {
		panic(fmt.Sprintf("invalid plugin config: %+v", *section))
	}
	this.class = section.String("class", "")
	if this.class == "" {
		this.class = this.name
	}
	this.enabled = section.Bool("enabled", true)
}

// A helper object to support delayed plugin creation.
type pluginWrapper struct {
	name          string
	configCreator func() *conf.Conf
	pluginCreator func() Plugin
}

func (this *Engine) LoadPlugins(section *conf.Conf) {
	pluginCommon := new(pluginCommon)
	pluginCommon.load(section)

	wrapper := new(pluginWrapper)
	var present bool
	if wrapper.pluginCreator, present = registeredPlugins[pluginCommon.class]; !present {
		panic("unknown plugin type: " + pluginCommon.class)
	}
	wrapper.name = pluginCommon.name
	wrapper.configCreator = func() *conf.Conf { return section }

	plugin := wrapper.pluginCreator()
	plugin.Init(section)
}
