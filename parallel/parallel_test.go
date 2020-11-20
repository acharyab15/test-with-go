package parallel

import (
	"fmt"
	"testing"
)

// func TestSomething(t *testing.T) {
// 	t.Parallel()
// 	time.Sleep(time.Second)
// 	// t.Fatal("not implemented")
// }

// func TestA(t *testing.T) {
// 	t.Parallel()
// 	time.Sleep(time.Second)
// }

// func TestB(t *testing.T) {
// 	fmt.Println("setup")
// 	defer fmt.Println("deferred teardown")
// 	t.Run("group", func(t *testing.T) {
// 		t.Run("sub1", func(t *testing.T) {
// 			t.Parallel()
// 			// run sub1
// 			time.Sleep(time.Second)
// 			fmt.Println("sub1 done")
// 		})
// 		t.Run("sub2", func(t *testing.T) {
// 			t.Parallel()
// 			// run sub2
// 			time.Sleep(time.Second)
// 			fmt.Println("sub2 done")
// 		})
// 	})
// 	fmt.Println("teardown")
// }

func TestGotcha(t *testing.T) {
	testCases := []struct {
		arg  int
		want int
	}{
		{2, 4},
		{3, 9},
		{4, 16},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("i=%d", tc.arg), func(t *testing.T) {
			c := tc // copy value for parallel tests; call before calling parallel
			t.Parallel()
			t.Logf("Testing with: arg=%d, want=%d", c.arg, c.want)
			if c.arg*c.arg != c.want {
				t.Errorf("%d^2 != %d", c.arg, c.want)
			}
			fmt.Println(c.arg, c.want)
		})
	}
}
