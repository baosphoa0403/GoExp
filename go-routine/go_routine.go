package goroutine

import (
	"fmt"
	"sync"
)

func SafeRunGoRoutine(callback func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("😱 Recovered from panic:", r)
			}
		}()

		callback()
	}()
}

func SafeRunGoRoutineWait(wg *sync.WaitGroup, callback func()) {
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("😱 Recovered from panic:", r)
			}
			wg.Done()
		}()

		callback()
	}()
}

func sayHello() {
	defer panic("lỗi nè")

	fmt.Println("Hello")
}

func Goroutine() {
	var wg sync.WaitGroup
	fmt.Printf("before sayHello\n")

	// will run background Not Wait
	// SafeRunGoRoutine(sayHello)

	// will run background Wait

	SafeRunGoRoutineWait(&wg, sayHello)

	fmt.Printf("after sayHello\n")
	wg.Wait()
}
