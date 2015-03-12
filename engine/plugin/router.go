package plugin

type router struct {
	input chan *PipelinePack
}

func newRouter(sz int) *router {
	return &router{input: make(chan *PipelinePack, sz)}
}

func (this *router) Run() {

LOOP:
	for {
		select {
		case pack, ok := <-this.input:
			if !ok {
				break LOOP
			}

			// find matchers and feed its InChan with pack
			m := matcher{}
			if m.Match(pack) {
				pack.IncRef()
				m.r.InChan() <- pack
			}

			pack.Recycle()
		}
	}

}
