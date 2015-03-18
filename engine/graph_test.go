package engine

import (
	"testing"
)

func TestPercentile(t *testing.T) {
	data := []uint64{4, 23, 32, 33, 88, 100}
	t.Logf("33  %+v", percentile(33., data))
	t.Logf("80  %+v", percentile(80., data))
	t.Logf("90  %+v", percentile(90., data))
	t.Logf("100 %+v", percentile(100., data))
}
