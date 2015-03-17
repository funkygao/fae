package plugin

import (
	"fmt"
	conf "github.com/funkygao/jsconf"
)

type Plugin interface {
	Init(section *conf.Conf)
	Run(r PluginRunner) error
}

func RegisterPlugin(name string, factory func() Plugin) {
	if _, present := registeredPlugins[name]; present {
		panic(fmt.Sprintf("plugin[%s] cannot register twice", name))
	}

	registeredPlugins[name] = factory
}

func LoadPlugins(cf *conf.Conf) {
	for i := 0; i < PackRecyclePoolSize; i++ {
		pack := NewPipelinePack(packRecycleChan)
		packRecycleChan <- pack
	}
	for i := 0; i < len(cf.List("plugins", nil)); i++ {
		section, err := cf.Section(fmt.Sprintf("plugins[%d]", i))
		if err != nil {
			panic(err)
		}

		name := section.String("name", "")
		if name == "" {
			panic("empty plugin name")
		}

		loadOnePlugin(name, section)
	}
}

func loadOnePlugin(name string, cf *conf.Conf) {
	wrapper := new(pluginWrapper)
	var present bool
	if wrapper.pluginCreator, present = registeredPlugins[name]; !present {
		panic("unknown plugin type: " + name)
	}
	wrapper.name = name
	wrapper.configCreator = func() *conf.Conf {
		return cf
	}
	wrappers[name] = wrapper

	plugin := wrapper.pluginCreator()
	plugin.Init(cf)

	runners[name] = newRunner(name, plugin)
}
