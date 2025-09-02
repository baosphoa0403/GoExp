package go101

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ctxKey string

const workerKey ctxKey = "worker-index"

func DoTaskCtx(ctx context.Context, url string) {
	start := time.Now()
	value := ctx.Value(workerKey)
	fmt.Println("get value from ctx value: ", value)
	// Tạo request có context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("req error:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Nếu timeout, lỗi sẽ là "context deadline exceeded"
		fmt.Printf("[timeout/error] %s (%v)\n", url, err)
		return
	}
	defer resp.Body.Close()

	_, _ = io.ReadAll(resp.Body)

	elapsed := time.Since(start)
	fmt.Printf("[ok] %s took %v\n", url, elapsed)
}

func WorkerCtx(i int, wg *sync.WaitGroup, ch <-chan string) {

	fmt.Println("Worker start index: ", i)
	for v := range ch {
		ctxValue := context.WithValue(context.Background(), workerKey, fmt.Sprintf("%d", i))
		ctx, cancel := context.WithTimeout(ctxValue, time.Second*1)
		DoTaskCtx(ctx, v)
		cancel() // giải phóng tài nguyên
	}

	defer wg.Done()
}

func HandleTaskWithContextTimeout() {
	var wg sync.WaitGroup
	workerNum := 2
	ch := make(chan string)
	urls := []string{"https://www.google.com/?hl=vi", "https://www.youtube.com/watch?v=zcNSD1QoBEg&list=RDET8RRB4srEU&index=2", "https://httpbin.org/delay/3"}

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			WorkerCtx(i, wg, ch)
		}(i, &wg)
	}

	for _, v := range urls {
		ch <- v
	}
	close(ch)

	wg.Wait()
}
