package blackrock

import (
	"testing"
)

func TestBlackrockWorks(t *testing.T) {
	br := Init(16777216, 45, 4)
	NumbersDelivered := make(map[int]int)

	for i := 0; i < 16777216; i++ {
		res := int(br.Shuffle(i))
		NumbersDelivered[res]++
	}

	for k, v := range NumbersDelivered {
		if v != 1 {
			t.Fatalf("The number %v came up twice in RNG", k)
		}
	}

	for i := 0; i < 16777216; i++ {
		if NumbersDelivered[i] == 0 {
			t.Fatalf("The number %v did not come up", i)
		}
	}
}
