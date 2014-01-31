package engine

import (
//conf "github.com/daviddengcn/go-ljson-conf"
)

type Engine struct {
	conf *Config
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	return
}
