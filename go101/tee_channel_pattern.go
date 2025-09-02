package go101

import (
	"fmt"
)

func nunberStreamStage1() <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	return ch
}

func teeChannel(in <-chan int) (<-chan int, <-chan int) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		defer close(ch2)

		for v := range in {
			ch1 <- v
			ch2 <- v
		}

	}()

	return ch1, ch2
}

func HandleTeeChannel() {
	out := nunberStreamStage1()

	ch1, ch2 := teeChannel(out)

	total := 0

	for {
		value1, ok1 := <-ch1
		value2, ok2 := <-ch2
		if !ok1 || !ok2 {
			break
		}

		fmt.Println("value1: ", value1)
		total += value2
		fmt.Println("total: ", total)
	}
}
