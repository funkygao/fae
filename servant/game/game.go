package game

import (
	"github.com/funkygao/fae/config"
	"github.com/funkygao/fae/servant/redis"
	"github.com/funkygao/metrics"
)

type Game struct {
	nameGen      *NameGen
	NameDbLoaded bool

	lock     *Lock
	register *Register
	presence *Presence

	phpLatency     metrics.Histogram // in ms
	phpPayloadSize metrics.Histogram // in bytes
}

func New(cf *config.ConfigGame, redis *redis.Client) *Game {
	this := new(Game)
	this.nameGen = newNameGen(cf.NamegenLength)
	this.lock = newLock(cf)
	this.register = newRegister(cf, redis)
	this.presence = newPresence()

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

func (this *Game) Register(typ string) (int64, error) {
	return this.register.Register(typ)
}

func (this *Game) CheckIn(uid int64) {
	this.presence.Update(uid)
}

func (this *Game) OnlineStatus(uids []int64) []bool {
	return this.presence.Onlines(uids)
}

func (this *Game) PhpPayloadHistogram() metrics.Histogram {
	return this.phpPayloadSize
}

func (this *Game) PhpLatencyHistogram() metrics.Histogram {
	return this.phpLatency
}
