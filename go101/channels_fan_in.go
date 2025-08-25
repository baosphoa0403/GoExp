package go101

import (
	"fmt"
	"sync"
)

func WorkerFanIn(wg *sync.WaitGroup, index int, out chan<- int) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		if index%2 == 0 && i%2 == 0 {
			out <- i
		} else {
			out <- i
		}
	}
}

func FanIn() {
	out := make(chan int)
	var wg sync.WaitGroup

	workerNum := 2

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go WorkerFanIn(&wg, i, out)
	}

	go func() {
		fmt.Println("zoo final")
		wg.Wait()
		close(out)
	}()

	for {
		value, ok := <-out
		if !ok {
			return
		}
		fmt.Println("value: ", value)
	}

}
