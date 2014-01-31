package engine

type Engine struct {
	conf *Config
}

func NewEngine() (this *Engine) {
	this = new(Engine)
	return
}
