package go101

import (
	"fmt"
	"sync"
)

func TeeMultiple[T any](in <-chan T, n int) []<-chan T {
	outs := make([]chan T, n)
	resChs := make([]<-chan T, n)

	for i := 0; i < n; i++ {
		outs[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, ch := range outs {
				close(ch)
			}
		}()

		for v := range in {
			for i := 0; i < n; i++ {
				outs[i] <- v
			}
		}
	}()

	for i := range outs {
		resChs[i] = outs[i]
	}

	return resChs
}

func numberStreamStage1() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			ch <- i
		}
	}()
	return ch
}

func HandleMultiTee() {
	in := numberStreamStage1()
	outs := TeeMultiple(in, 3) // tạo 3 channel output
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		for v := range outs[0] {
			fmt.Println("consumer1 got", v)
		}
		wg.Done()
	}()

	go func() {
		for v := range outs[1] {
			fmt.Println("consumer2 got", v)
		}
		wg.Done()
	}()

	go func() {
		for v := range outs[2] {
			fmt.Println("consumer3 got", v)
		}
		wg.Done()
	}()

	wg.Wait()
}
