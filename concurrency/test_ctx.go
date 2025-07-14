package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func HandleTask(c chan string, tasks []int, ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < len(tasks); i++ {
		wg.Add(1)
		go func(i int, ctx context.Context) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			ctxChild, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			if i == 3 {
				select {
				case <-time.After(2 * time.Second):
					c <- fmt.Sprintf("Done task %d", i)
					fmt.Println("✅ Task completed after sleep at index:", i)
				case <-ctxChild.Done():
					fmt.Println("❌ Task canceled (timeout) at index:", i, "-", ctxChild.Err())
					return
				}
			}

			if i == 1 {
				panic(fmt.Sprintf("error at: %d", i))
			}

			c <- fmt.Sprintf("Done task %d", i)
		}(i, ctx)
	}

	wg.Wait()
	fmt.Println("close channel")
	close(c)
}

func TestCtx() {
	c := make(chan string, 10)
	parentCtx, cancel := context.WithCancel(context.Background())
	tasks := []int{1, 2, 3, 4, 5}

	go HandleTask(c, tasks, parentCtx)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("🔔 Manually cancel context after 2s")
		cancel()
	}()

	for {
		select {
		case <-parentCtx.Done():
			fmt.Println("⚠️ Main context canceled:", parentCtx.Err())
			fmt.Println("✅ Done main goroutine")
			return
		case res, ok := <-c:
			if !ok {
				fmt.Println("End task res: ", res, "- ok: ", ok)
				return
			}
			fmt.Println("res: ", res, "- ok: ", ok)
		}
	}
}
