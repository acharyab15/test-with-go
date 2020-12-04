package draw_test

import (
	"fmt"
	"image/draw"
	"testing"

	twgdraw "github.com/acharyab15/test-with-go/draw"
)

func TestFibGradient(t *testing.T) {
	var im draw.Image
	twgdraw.FibGradient(im)
}

func TestFibFunc(t *testing.T) {
	fmt.Println(twgdraw.A)
	got := twgdraw.Fib(2)
	if got != 1 {
		t.Errorf("Fib(2) = %d, want 1", got)
	}

}

// func TestInfo(t *testing.T) {
// 	d := twgdraw.Dog{}
// 	twgdraw.Info(d)
// }
