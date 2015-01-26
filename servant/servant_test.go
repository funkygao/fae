package servant

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestGolangStringCmp(t *testing.T) {
	assert.Equal(t, true, "1" == "1")
	assert.Equal(t, true, "2" > "1")
	assert.Equal(t, false, "1" > "1")
	assert.Equal(t, true, "22" > "2")
}
