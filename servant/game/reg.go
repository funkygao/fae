package game

import (
	log "github.com/funkygao/log4go"
)

type Register struct {
	maxPerKingdom int
	tilesInK      int
	tiles         map[int]map[int]bool
	k             int
}

func newRegister(maxPerKingdom int) *Register {
	this := new(Register)
	this.maxPerKingdom = maxPerKingdom
	this.tiles = make(map[int]map[int]bool)
	this.loadSnapshot()
	return this
}

func (this *Register) loadSnapshot() {
	// fill k, tiles, tilesInK
	log.Debug("register snapshot loaded")
}

func (this *Register) RegTile() (k, x, y int) {
	this.tilesInK++
	if this.tilesInK >= this.maxPerKingdom {
		this.k++
		this.tilesInK = 0
	}

	k = this.k

	return
}
