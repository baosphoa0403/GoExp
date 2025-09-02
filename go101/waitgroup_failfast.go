package go101

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func PromiseAll(callbacks ...func(ctx context.Context, workerId int) <-chan Result) ([]Result, error) {
	results := make([]Result, len(callbacks))
	var wg sync.WaitGroup
	cancels := make([]context.CancelFunc, len(callbacks))
	errCh := make(chan error, 1)
	for i := 0; i < len(callbacks); i++ {
		wg.Add(1)
		ctx, cancel := context.WithCancel(context.Background())
		cancels[i] = cancel
		go func(ctx context.Context, i int) {
			defer wg.Done()
			data := <-callbacks[i](ctx, i)
			if data.error != nil {
				for _, cancel := range cancels {
					if cancel != nil {
						cancel()
					}
				}
				select {
				case errCh <- data.error:
				default:
				}
				return
			}
			results[i] = data
		}(ctx, i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-done:
		fmt.Println()
		return results, nil
	}
}

func CallApi(ctx context.Context, workerId int, taskName string, err error) <-chan Result {
	out := make(chan Result, 1)
	go func() {
		defer close(out)

		if err != nil {
			out <- Result{
				data:  nil,
				error: fmt.Errorf("error worker %d - error: %s", workerId, err.Error()),
			}
			return
		}

		select {
		case <-ctx.Done():
			fmt.Println("context cancel")
			out <- Result{
				data:  nil,
				error: ctx.Err(),
			}
		case <-time.After(4 * time.Second):
			fmt.Printf("Worker %d running task: %s...\n", workerId, taskName)
			out <- Result{
				data:  fmt.Sprintf("run success at workerId: %d", workerId),
				error: nil,
			}
		}
	}()
	return out

}

func HandlePromiseAllGolang() {

	res, err := PromiseAll(
		func(ctx context.Context, workerId int) <-chan Result {
			return CallApi(ctx, workerId, "task 1", nil)
		},
		func(ctx context.Context, workerId int) <-chan Result {
			return CallApi(ctx, workerId, "task 2", nil)
		})

	fmt.Println("res: ", res)
	if err != nil {
		fmt.Println("err: ", err)
	}

}
