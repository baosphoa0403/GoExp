package go101

import (
	"fmt"
	"sync"
	"time"
)

var bufferSize = 3

func Stage4() <-chan int {
	// Stage 1: sinh ra 100 số nguyên từ 1 → 100.
	ch := make(chan int, bufferSize)

	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func Stage5(in <-chan int) <-chan int {
	// Stage 2: lọc ra số chẵn.
	ch := make(chan int, bufferSize)
	var wg sync.WaitGroup

	workerNum := 3

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				if v%2 == 0 {
					time.Sleep(time.Second * 2)
					ch <- v
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func Stage6(in <-chan int) {
	// Stage 3: bình phương số đó và in ra.

	for v := range in {
		fmt.Println("final stage 6 receive value: ", v*v)

	}

}

func PieplineBuffered() {
	ch := Stage4()
	ch = Stage5(ch)
	Stage6(ch)
}
