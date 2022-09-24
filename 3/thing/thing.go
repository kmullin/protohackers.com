package thing

import (
	"fmt"
	"sync"
)

func Merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			out <- v
		}
		fmt.Println("stopping")
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// close output channel once all producers are done
	go func() {
		wg.Wait()
		fmt.Println("closing output chan")
		close(out)
	}()
	return out
}
