package go101

import (
	"fmt"
	"sync"
	"time"
)

func HandleSemaphore(limit int) {
	start := time.Now()

	var wg sync.WaitGroup
	sem := make(chan struct{}, limit) // semaphore với capacity = limit

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func(taskID int) {
			defer wg.Done()	

			// acquire slot
			sem <- struct{}{}

			fmt.Println("task started:", taskID)
			time.Sleep(time.Second * 1) // mô phỏng công việc nặng
			fmt.Println("task finished:", taskID)

			// release slot
			<-sem
		}(i)
	}

	wg.Wait()
	fmt.Printf("All tasks done with limit=%d, elapsed=%v\n", limit, time.Since(start))
}
