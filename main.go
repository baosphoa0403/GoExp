package main

import (
	"context"
	"fmt"
	"time"

	"giabao.com/utils"
)

func main() {
	// goroutine.Goroutine()
	// pdf()
	// lock.MutexLock()

	// for i := 0; i < 15; i++ {
	// 	jobID := i
	// 	pool.Submit(func() {
	// 		fmt.Printf("Job %d running\n", jobID)
	// 		time.Sleep(500 * time.Millisecond)

	// 		if jobID == 7 {
	// 			panic("Job 7 bị lỗi nè")
	// 		}

	// 		fmt.Printf("Job %d done\n", jobID)
	// 	})
	// }

	// Submit job có context timeout

	// fastJob()
	// longerJob()
	// queueFullCase()
	jobLostAfterStop()
}

func fastJob() {
	pool := utils.NewPool(10, 100)
	for i := 0; i < 100; i++ {
		index := i
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := pool.SubmitWithContext(ctx, func(ctx context.Context) {
			select {
			case <-time.After(10 * time.Millisecond): // giả lập job chạy lâu
				fmt.Println("✅ Job context xong:", index)
			case <-ctx.Done():
				fmt.Println("❌ Job context HUỶ:", index)
			}
		})
		if err != nil {
			fmt.Println("SubmitWithContext error:", err)
		}
	}

	pool.Stop()
	fmt.Println("🎉 Done Case 1")
}

func longerJob() {
	pool := utils.NewPool(10, 100)
	for i := 0; i < 100; i++ {
		index := i
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := pool.SubmitWithContext(ctx, func(ctx context.Context) {
			select {
			case <-time.After(500 * time.Millisecond): // giả lập job chạy lâu
				fmt.Println("✅ Job context xong:", index)
			case <-ctx.Done():
				fmt.Println("❌ Job context HUỶ:", index)
			}
		})
		if err != nil {
			fmt.Println("SubmitWithContext error:", err)
		}
	}

	pool.Stop()
	fmt.Println("🎉 Done Case 1")
}

func queueFullCase() {
	pool := utils.NewPool(1, 1) // chỉ 1 worker + queue size rất nhỏ

	ctx := context.Background()

	// Job đầu tiên → chiếm worker
	pool.SubmitWithContext(ctx, func(ctx context.Context) {
		time.Sleep(1 * time.Second)
	})

	// Job thứ hai → chiếm queue
	pool.SubmitWithContext(ctx, func(ctx context.Context) {
		time.Sleep(1 * time.Second)
	})

	// Job thứ ba → vượt queue
	err := pool.SubmitWithContext(ctx, func(ctx context.Context) {
		fmt.Println("❌ This should NOT be printed")
	})
	if err != nil {
		fmt.Println("✅ Case 3 caught:", err) // should print job queue full
	}

	pool.Stop()
}

func jobLostAfterStop() {
	pool := utils.NewPool(1, 10)

	for i := 0; i < 5; i++ {
		index := i
		ctx := context.Background()
		err := pool.SubmitWithContext(ctx, func(ctx context.Context) {
			time.Sleep(20 * time.Millisecond)
			fmt.Println("👀 Job chạy:", index)
		})
		if err != nil {
			fmt.Println("⚠ Lỗi submit:", err)
		}
	}

	time.Sleep(5 * time.Millisecond)
	pool.Stop()
	fmt.Println("🎉 Done Case 4")
}
