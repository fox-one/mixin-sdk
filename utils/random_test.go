package utils

import (
	"testing"
)

func TestRandInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		min := RandInt(0, 10)
		max := RandInt(0, 10)

		if min >= max {
			min, max = max, min+1
		}

		rand := RandInt(min, max)

		if rand < min || rand >= max {
			t.Error("rand int out of bounds")
		}
	}
}
