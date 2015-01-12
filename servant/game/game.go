package game

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/metrics"
)

type Game struct {
	nameGen      *NameGen
	NameDbLoaded bool

	lock     *Lock
	register *Register

	phpLatency     metrics.Histogram // in ms
	phpPayloadSize metrics.Histogram // in bytes
}

func New(nameSlot int, lk *config.ConfigLock) *Game {
	this := new(Game)
	this.nameGen = newNameGen(nameSlot)
	this.lock = newLock(lk)
	this.register = newRegister()

	this.phpLatency = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("php.latency", this.phpLatency)
	this.phpPayloadSize = metrics.NewHistogram(
		metrics.NewExpDecaySample(1028, 0.015))
	metrics.Register("php.payload", this.phpPayloadSize)

	return this
}

func (this *Game) NextName() string {
	return this.nameGen.Next()
}

func (this *Game) SetNameBusy(name string) error {
	return this.nameGen.SetBusy(name)
}

func (this *Game) Lock(key string) (success bool) {
	return this.lock.Lock(key)
}

func (this *Game) Unlock(key string) {
	this.lock.Unlock(key)
}

func (this *Game) UpdatePhpLatency(latency int64) {
	this.phpLatency.Update(latency)
}

func (this *Game) UpdatePhpPayloadSize(bytes int64) {
	this.phpPayloadSize.Update(bytes)
}

func (this *Game) RegTile() (k, x, y int) {
	return this.register.RegTile()
}
