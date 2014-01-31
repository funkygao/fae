package engine

import (
//conf "github.com/daviddengcn/go-ljson-conf"
)

type Engine struct {
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	return
}

func (this *Engine) LoadConfigFile(fn string) *Engine {
	return this
}

func (this *Engine) ServeForever() {

}
