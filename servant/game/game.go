package game

type Game struct {
	nameGen      *NameGen
	NameDbLoaded bool
}

func New(nameSlot int) *Game {
	this := new(Game)
	this.nameGen = newNameGen(nameSlot)
	return this
}

func (this *Game) NextName() string {
	return this.nameGen.Next()
}

func (this *Game) SetNameBusy(name string) error {
	return this.nameGen.SetBusy(name)
}
