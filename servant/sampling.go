package servant

import (
	"math/rand"
)

func sampleRateSatisfied(rate int) bool {
	return rand.Intn(100) <= rate
}
