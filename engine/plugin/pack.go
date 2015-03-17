package plugin

import (
	"sync/atomic"
)

type PipelinePack struct {
	// Where to put back myself when reference count zeros
	RecycleChan chan *PipelinePack

	// Reference counter, internal GC
	RefCount int32

	Ident string
	Data  interface{}
}

func NewPipelinePack(recyleChan chan *PipelinePack) *PipelinePack {
	return &PipelinePack{
		RecycleChan: recyleChan,
		RefCount:    int32(1), // if 0, will be GC'ed
	}
}

func (this *PipelinePack) IncRef() {
	atomic.AddInt32(&this.RefCount, 1)
}

func (this *PipelinePack) Recycle() {
	rc := atomic.AddInt32(&this.RefCount, -1)
	if rc == 0 {
		this.Reset()

		this.RecycleChan <- this
	} else if rc < 0 {
		// panic? TODO
	}
}

func (this *PipelinePack) Reset() {
	atomic.StoreInt32(&this.RefCount, 1)
	this.Data = nil
	this.Ident = ""
}
