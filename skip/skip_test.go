package skip

import "testing"

// var shouldBeSkipped = true

func TestThing(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Log("this test ran!")
}
