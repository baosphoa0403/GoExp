package go101

import (
	"fmt"
	"sync"
	"time"
)

// Worker: đọc job từ channel tới khi channel đóng
func Worker(wg *sync.WaitGroup, index int, jobs <-chan string) {
	fmt.Printf("worker %d start\n", index)
	for job := range jobs {
		// giả lập xử lý job
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("worker %d processing %s\n", index, job)
		wg.Done()
	}
	fmt.Printf("worker %d exit (channel closed)\n", index)
}

// ================= Version 1 =================
// Producer chạy trong main goroutine
func FanOutMainProducer() {
	var wg sync.WaitGroup
	jobs := make(chan string)
	workerNum := 3

	// start workers
	for i := 0; i < workerNum; i++ {
		go Worker(&wg, i, jobs)
	}

	// producer trong main
	for i := 0; i < 10; i++ {
		wg.Add(1) // Add ngay trước khi gửi job
		jobs <- fmt.Sprintf("[main producer] job %d", i)
	}
	close(jobs)

	wg.Wait()
	fmt.Println("FanOutMainProducer done")
}

// ================= Version 2 =================
// Producer chạy trong goroutine riêng
func FanOutGoroutineProducer() {
	var wg sync.WaitGroup
	jobs := make(chan string)
	workerNum := 3
	jobCount := 10

	// start workers
	for i := 0; i < workerNum; i++ {
		go Worker(&wg, i, jobs)
	}

	// Add đủ trước, vì producer chạy trong goroutine riêng
	wg.Add(jobCount)

	// producer tách goroutine
	go func() {
		for i := 0; i < jobCount; i++ {
			jobs <- fmt.Sprintf("[goroutine producer] job %d", i)
		}
		close(jobs)
	}()

	wg.Wait()
	fmt.Println("FanOutGoroutineProducer done")
}

// Producer trong main → dùng khi biết trước số lượng job, finite.

// Producer goroutine riêng → khi job dynamic hoặc infinite stream.

// WaitGroup → chỉ phù hợp khi số job finite.

// Context/Signal → dùng để quản lý khi job vô hạn.
