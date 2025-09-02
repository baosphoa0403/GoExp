package go101

import (
	"context"
	"fmt"
	"time"
)

// func TaskCtxCancel(ctx context.Context, wg *sync.WaitGroup) {
// 	for i := 0; ; i++ {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("done: ", ctx.Err().Error())
// 			wg.Done()
// 			return
// 		default:
// 			fmt.Println("Ping:", i)            // chỉ in khi chưa cancel
// 			time.Sleep(200 * time.Millisecond) // tránh spam CPU
// 		}
// 	}
// }

// func HandleTaskCancel() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	var wg sync.WaitGroup

// 	wg.Add(1)
// 	go TaskCtxCancel(ctx, &wg)

// 	go func() {
// 		fmt.Println("goroutine cancel run")
// 		time.Sleep(time.Second * 2)
// 		fmt.Println("call cancel after 3 second")
// 		cancel()
// 	}()

// 	wg.Wait()

// }

func TaskCtxCancel(ctx context.Context) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("done: ", ctx.Err().Error())
			return
		default:
			fmt.Println("Ping:", i)            // chỉ in khi chưa cancel
			time.Sleep(200 * time.Millisecond) // tránh spam CPU
		}
	}
}

func HandleTaskCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	go TaskCtxCancel(ctx)

	select {
	case <-time.After(time.Second * 1):
		fmt.Println("timeout after 1s")
		cancel()
		return
	}

}
