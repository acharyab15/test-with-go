package random

import (
	"math/rand"
	"testing"
)

func TestPick(t *testing.T) {
	// seed := time.Now().UnixNano()
	// printing out/ keeping track of a seed could be helpful for repeating tests
	var seed int64 = 1605885692917778000
	r := rand.New(rand.NewSource(seed))
	arg := make([]int, 25)
	for i := 0; i < 25; i++ {
		arg[i] = r.Int()
	}
	got := Pick(arg)
	for _, v := range arg {
		if got == v {
			return
		}
	}
	t.Errorf("Pick(seed=%d)=%d; not in slice", seed, got)
}
