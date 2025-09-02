package go101

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID int
}

// Worker function
func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// mô phỏng xử lý mất thời gian ngẫu nhiên
		fmt.Printf("Worker %d bắt đầu xử lý job %d\n", id, job.ID)
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		fmt.Printf("Worker %d xong job %d\n", id, job.ID)
	}
}

// Bounded Work Queue
func boundedWorkQueue(numWorkers int, queueSize int, jobsToProduce int) {
	jobs := make(chan Job, queueSize) // queue có capacity
	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Producer
	for j := 1; j <= jobsToProduce; j++ {
		job := Job{ID: j}
		fmt.Printf("Producer gửi job %d\n", j)
		jobs <- job // block nếu queue đầy
	}

	close(jobs) // đóng queue
	wg.Wait()
}

func HandleBoundedWorkQueue() {
	boundedWorkQueue(3, 5, 15)
}
