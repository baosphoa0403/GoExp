package go101

import (
	"fmt"
	"sync"
	"time"
)

var bufferSize = 10_000_000

func Stage4() <-chan int {
	// Stage 1: sinh ra 100 số nguyên từ 1 → 100.
	ch := make(chan int, bufferSize)
	count := 0
	go func() {
		for i := 1; ; i++ {
			if count == bufferSize {
				fmt.Println("sleep max count: ", count)
				time.Sleep(time.Second * 2)
				count = 0
			}
			ch <- i
			count += 1
		}
		// close(ch)
	}()

	return ch
}

func Stage5(in <-chan int) <-chan int {
	// Stage 2: lọc ra số chẵn.
	ch := make(chan int, bufferSize)
	var wg sync.WaitGroup

	workerNum := 5

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				if v%2 == 0 {
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
	// runtime.GOMAXPROCS(1) // chỉ dùng 2 core
	ch := Stage4()
	ch = Stage5(ch)
	Stage6(ch)
}
