package thing

import (
	"fmt"
	"sync"
)

type ChanReader []<-chan int

func (c ChanReader) Something() <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(c))
	for _, cr := range c {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
			fmt.Println("stopping")
		}(cr)
	}
	go func() {
		wg.Wait()
		fmt.Println("closing output chan")
		close(out)
	}()
	return out
}
