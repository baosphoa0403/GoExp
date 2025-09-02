package go101

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SettledResult struct {
	Status string
	Value  any
	Error  error
}

func PromiseAllSettled(callbacks ...func(ctx context.Context, workerId int) <-chan Result) []SettledResult {
	results := make([]SettledResult, len(callbacks))
	var wg sync.WaitGroup

	for i := 0; i < len(callbacks); i++ {
		wg.Add(1)
		ctx := context.Background() // allSettled không cancel
		go func(i int) {
			defer wg.Done()
			res := <-callbacks[i](ctx, i)
			if res.error != nil {
				results[i] = SettledResult{
					Status: "rejected",
					Error:  res.error,
				}
			} else {
				results[i] = SettledResult{
					Status: "fulfilled",
					Value:  res.data,
				}
			}
		}(i)
	}

	wg.Wait()
	return results
}

func CallApi1(ctx context.Context, workerId int, taskName string, err error) <-chan Result {
	out := make(chan Result, 1)
	go func() {
		defer close(out)
		fmt.Printf("Worker %d running task: %s...\n", workerId, taskName)

		if err != nil {
			out <- Result{
				data:  nil,
				error: fmt.Errorf("error worker %d - %s", workerId, err),
			}
			return
		}

		time.Sleep(2 * time.Second)
		out <- Result{
			data:  fmt.Sprintf("run success at workerId: %d", workerId),
			error: nil,
		}
	}()
	return out
}

func HandlePromiseAllSettled() {
	results := PromiseAllSettled(
		func(ctx context.Context, workerId int) <-chan Result {
			return CallApi1(ctx, workerId, "task 1", nil) // success
		},
		func(ctx context.Context, workerId int) <-chan Result {
			return CallApi1(ctx, workerId, "task 2", fmt.Errorf("fail nè")) // fail
		},
	)

	fmt.Println("results: ", results)

	for i, r := range results {
		if r.Status == "fulfilled" {
			fmt.Printf("Task %d success: %v\n", i, r.Value)
		} else {
			fmt.Printf("Task %d failed: %v\n", i, r.Error)
		}
	}
}
