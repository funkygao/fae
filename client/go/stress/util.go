package main

import (
	"math/rand"
)

func sampling(rate int) bool {
	return rand.Intn(rate) == 1
}
