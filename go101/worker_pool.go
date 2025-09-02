package go101

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func DoTask(url string, wg1 *sync.WaitGroup) {
	defer wg1.Done()

	start := time.Now() // bắt đầu đo thời gian

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// đọc nội dung (optional, có thể bỏ nếu chỉ cần đo thời gian)
	_, _ = io.ReadAll(resp.Body)

	elapsed := time.Since(start) // tính thời gian request
	fmt.Printf("Request to %s took %v\n", url, elapsed)
}

func DoTask1(url string) {

	start := time.Now() // bắt đầu đo thời gian

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// đọc nội dung (optional, có thể bỏ nếu chỉ cần đo thời gian)
	_, _ = io.ReadAll(resp.Body)

	elapsed := time.Since(start) // tính thời gian request
	fmt.Printf("Request to %s took %v\n", url, elapsed)
}

func WorkerInPool(index int, wg *sync.WaitGroup, ch <-chan string) {
	var wg1 sync.WaitGroup
	fmt.Println("worker done: ", index)
	for v := range ch {
		wg1.Add(1)
		go func() {
			DoTask(v, &wg1)
		}()
	}

	wg1.Wait()
	wg.Done()
}

func WorkerInPool1(index int, wg *sync.WaitGroup, ch <-chan string) {
	fmt.Println("worker done: ", index)
	for v := range ch {
		DoTask1(v)
	}

	wg.Done()
}

func HandleWorkerPool() {
	workerNum := 3
	var wg sync.WaitGroup
	ch := make(chan string)
	urls := []string{"https://www.google.com/?hl=vi", "https://www.youtube.com/watch?v=zcNSD1QoBEg&list=RDET8RRB4srEU&index=2", "https://dantri.com.vn/", "https://tinhte.vn/"}

	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup, ch <-chan string) {
			WorkerInPool1(i, wg, ch)
		}(i, &wg, ch)
	}

	for _, v := range urls {
		ch <- v
	}
	close(ch)

	wg.Wait()

}
