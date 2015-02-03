package game

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestPresence(t *testing.T) {
	p := newPresence()
	p.Update(12)
	p.Update(3)
	assert.Equal(t, []bool{true, true}, p.Onlines([]int64{12, 3}))
	assert.Equal(t, []bool{true, false}, p.Onlines([]int64{12, 31}))
	assert.Equal(t, []bool{false, false}, p.Onlines([]int64{122, 23}))
}
