package namegen

import (
	"math/rand"
	"time"
)

// TODO space will not be used
const (
	NameCharMin uint8 = 33
	NameCharMax uint8 = 126
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type NameGen struct {
	bits [][]byte
}

func New(size int) (this *NameGen) {
	this = new(NameGen)
	this.bits = make([][]byte, size)
	for i, _ := range this.bits {
		this.bits[i] = make([]byte, (NameCharMax-NameCharMin)/8+1)

		for j := NameCharMin; j <= NameCharMax; j++ {
			this.bits[i][this.pos(j)] = 0
		}
	}

	return
}

func (this *NameGen) pos(c uint8) int {
	var x = c - NameCharMin

	return int(x / 8)

}

func (this *NameGen) Next() string {
	rv := ""
	for i := 0; i < len(this.bits); i++ {
		w := NameCharMin + uint8(rand.Int31n(int32(NameCharMax-NameCharMin)))

		rv += string(w)
	}
	return rv
}

func (this *NameGen) Size() int {
	return len(this.bits[0]) * len(this.bits)
}
