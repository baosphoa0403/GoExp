package concurrency

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func doSomeThing(index int, v string) (string, error) {
	fmt.Println("🚀 Start doSomeThing - index:", index, "value:", v)
	// if index == 2 {
	// 	return "", fmt.Errorf("loi ne at index %d", index)
	// }
	return fmt.Sprintf("done.%s", v), nil
}

func handleTaskGo(g *errgroup.Group, urls []string, ctx context.Context) {

	for index, v := range urls {
		g.Go(func() error {
			fmt.Println("zoo index: ", index)
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			if index == 1 {
				panic(fmt.Sprintf("error at: %d", index))
			}

			if index%2 == 0 {
				select {
				case <-ctx.Done():
					fmt.Println("⛔ Task canceled (even index)", index, ctx.Err())
					return ctx.Err()
				case <-time.After(5 * time.Second):
					fmt.Println("✅ Done after sleep index:", index)
					return nil
				}
			}

			select {
			case <-ctx.Done():
				fmt.Println("⛔ Task canceled due to context timeout", ctx.Err())
				return ctx.Err()
			default:
				result, err := doSomeThing(index, v)
				if err != nil {
					fmt.Println("❌ Task failed:", err)
					return err // ❗ propagate the error
				}
				fmt.Println("✅ Task success:", result)
				return nil
			}
		})
	}
}

func handleTaskTryGo(g *errgroup.Group, urls []string, ctx context.Context) {

	for index, v := range urls {
		ok := g.TryGo(func() error {
			fmt.Println("zoo index: ", index)
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			// if index == 1 {
			// 	panic(fmt.Sprintf("error at: %d", index))
			// }

			if index%2 == 0 {
				select {
				case <-ctx.Done():
					fmt.Println("⛔ Task canceled (even index)", index, ctx.Err())
					return ctx.Err()
				case <-time.After(5 * time.Second):
					fmt.Println("✅ Done after sleep index:", index)
					return nil
				}
			}

			select {
			case <-ctx.Done():
				fmt.Println("⛔ Task canceled due to context timeout", ctx.Err())
				return ctx.Err()
			default:
				result, err := doSomeThing(index, v)
				if err != nil {
					fmt.Println("❌ Task failed:", err)
					return err // ❗ propagate the error
				}
				fmt.Println("✅ Task success:", result)
				return nil
			}
		})

		fmt.Println("TryGo success?", ok) // 👉 false nếu context đã cancel
	}
}

func TestErrorGroup() {
	parentCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	g, ctx := errgroup.WithContext(parentCtx)

	g.SetLimit(1)
	urls := []string{"a.com", "b.com", "c.com", "d.com"}

	// handleTaskGo(g, urls, ctx)
	handleTaskTryGo(g, urls, ctx)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("⛔ Manually cancel context after 2s")
		cancel()
	}()

	if err := g.Wait(); err != nil {
		fmt.Println("❌ Error:", err)
	} else {
		fmt.Println("✅ All fetched successfully")
	}
}
