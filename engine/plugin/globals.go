package plugin

import (
	"regexp"
)

var (
	PackRecyclePoolSize = 200

	registeredPlugins = make(map[string]func() Plugin) // name:factory
	wrappers          = make(map[string]*pluginWrapper)
	runners           = make(map[string]PluginRunner)
	pluginTypeRegex   = regexp.MustCompile("^.*(Filter|Input|Output)$")
	packRecycleChan   = make(chan *PipelinePack, PackRecyclePoolSize)
	hub               = newRouter(PackRecyclePoolSize)
)
