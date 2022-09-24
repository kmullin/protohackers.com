package thing

import (
	"testing"
)

func TestThing(t *testing.T) {
	chanReader := make(ChanReader, 0)
	c1 := make(chan int)
	c2 := make(chan int)

	chanReader = append(chanReader, c1, c2)

	out := chanReader.Something()

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
