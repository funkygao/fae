package game

type Register struct {
}

func newRegister() *Register {
	this := new(Register)
	return this
}

func (this *Register) RegTile() (k, x, y int) {
	return
}
