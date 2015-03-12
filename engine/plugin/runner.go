package plugin

import (
	log "github.com/funkygao/log4go"
)

type PluginRunner interface {
	Name() string
	Plugin() Plugin
	InChan() chan *PipelinePack
	Inject(*PipelinePack)
	Run()
}

type runner struct {
	name   string
	plugin Plugin
	inChan chan *PipelinePack
}

func newRunner(name string, plugin Plugin) *runner {
	return &runner{
		name:   name,
		plugin: plugin,
		inChan: make(chan *PipelinePack, PackRecyclePoolSize),
	}
}

func (this *runner) Name() string {
	return this.name
}

func (this *runner) Plugin() Plugin {
	return this.plugin
}

func (this *runner) InChan() chan *PipelinePack {
	return this.inChan
}

func (this *runner) Inject(pack *PipelinePack) {
	hub.hub <- pack
}

func (this *runner) Run() {
	for {
		if err := this.plugin.Run(this); err != nil {
			log.Warn(err)
		}

		// restart this plugin

		// re-initialize my plugin using its wrapper
		wrapper := wrappers[this.name]
		this.plugin = wrapper.Create()
	}
}
