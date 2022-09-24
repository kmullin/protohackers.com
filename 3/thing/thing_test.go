package thing

import (
	"testing"
)

func TestMerge(t *testing.T) {
	c1 := make(chan int)
	c2 := make(chan int)
	out := Merge(c1, c2)

	c1 <- 10
	c2 <- 20

	num := <-out
	if num != 10 {
		t.Fatalf("num not 10")
	}
	num = <-out
	if num != 20 {
		t.Fatalf("num not 20")
	}
	close(c1)
	defer close(c2)

	c2 <- 30
	num = <-out
	if num != 30 {
		t.Fatalf("num not 30")
	}
}
