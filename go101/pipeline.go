package go101

import (
	"fmt"
	"sync"
)

func Stage1() <-chan int {
	// Stage 1: sinh ra 100 số nguyên từ 1 → 100.
	ch := make(chan int)

	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func Stage2(in <-chan int) <-chan int {
	// Stage 2: lọc ra số chẵn.
	ch := make(chan int)
	var wg sync.WaitGroup

	workerNum := 3

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("start worker")
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

func Stage3(in <-chan int) {
	// Stage 3: bình phương số đó và in ra.

	for v := range in {
		fmt.Println("value: ", v*v)
	}

}

func Piepline() {
	ch := Stage1()
	ch = Stage2(ch)
	Stage3(ch)
}
