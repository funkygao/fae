package plugin

type router struct {
	hub chan *PipelinePack
}

func newRouter(sz int) *router {
	return &router{hub: make(chan *PipelinePack, sz)}
}

func (this *router) Run() {

}
