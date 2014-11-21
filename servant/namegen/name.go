package namegen

import (
	"math/rand"
	"time"
)

// TODO space will not be used
const (
	NameCharMin = uint8('!') // 33, space is 32
	NameCharMax = uint8('~') // 126
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type NameGen struct {
	bits [][]byte
}

func New(slots int) (this *NameGen) {
	this = new(NameGen)
	this.bits = make([][]byte, slots)
	for i, _ := range this.bits {
		this.bits[i] = make([]byte, (NameCharMax-NameCharMin)/8+1)

		for char := NameCharMin; char <= NameCharMax; char++ {
			this.bits[i][this.pos(char)] = 0x00
		}
	}

	return
}

func (this *NameGen) offset(char uint8) uint {
	return uint(char - NameCharMin)
}

func (this *NameGen) pos(char uint8) uint {
	return uint(this.offset(char) / 8)
}

func (this *NameGen) index(char uint8) int {
	return int(this.offset(char) % 8)
}

func (this *NameGen) getBit(slot int, char uint8) uint8 {
	index, pos := this.index(char), this.pos(char)
	if slot > len(this.bits) || char < NameCharMin || char > NameCharMax {
		return 0
	}

	return (this.bits[slot][index] >> pos) & 0x01
}

func (this *NameGen) setBit(slot int, char uint8, value uint8) {
	index, pos := this.index(char), this.pos(char)
	if slot > len(this.bits) || char < NameCharMin || char > NameCharMax {
		return
	}

	// value can only be 0/1
	if value == 0 {
		this.bits[slot][index] &^= 0x01 << pos
	} else {
		this.bits[slot][index] |= 0x01 << pos
	}

}

func (this *NameGen) Next() string {
	rv := ""
	for i := 0; i < len(this.bits); i++ {
		w := NameCharMin + uint8(rand.Int31n(int32(NameCharMax-NameCharMin)))
		this.setBit(i, w, 1)

		rv += string(w)
	}
	return rv
}

func (this *NameGen) Size() int {
	return len(this.bits[0]) * len(this.bits) * 8
}
