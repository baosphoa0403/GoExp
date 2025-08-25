# 🕹️ Go Concurrency Guide

Tài liệu này tóm tắt các công cụ và pattern phổ biến để lập trình **concurrent** trong Go.  

---

## 1. Channels

### 📌 Use Cases
- **Fan-out**: một goroutine sinh dữ liệu → nhiều goroutines tiêu thụ.  
- **Fan-in**: nhiều goroutines sinh dữ liệu → gom về một channel duy nhất.  
- **Pipeline**: chuỗi goroutines, mỗi bước xử lý dữ liệu và gửi tiếp.  
- **Signal/Notify**: báo hiệu hoàn thành hoặc dừng.  

### 📌 Example (Fan-in)
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("worker %d processing job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= 5; a++ {
        fmt.Println(<-results)
    }
}
```

---

## 2. sync Package

### 🔒 Mutex
Bảo vệ biến dùng chung giữa nhiều goroutines.
```go
var mu sync.Mutex
var n int

func main() {
    for i := 0; i < 1000; i++ {
        go func() {
            mu.Lock()
            n++
            mu.Unlock()
        }()
    }
    time.Sleep(time.Second)
    fmt.Println("n =", n)
}
```

### 📚 RWMutex
- Nhiều readers có thể truy cập cùng lúc.  
- Writer phải độc quyền.  

### ⏳ WaitGroup
Chờ nhiều goroutines kết thúc.
```go
var wg sync.WaitGroup

func worker(id int) {
    defer wg.Done()
    fmt.Println("Worker", id, "started")
    time.Sleep(time.Second)
    fmt.Println("Worker", id, "done")
}

func main() {
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i)
    }
    wg.Wait()
    fmt.Println("All workers finished")
}
```

### 🟢 Once
Đảm bảo một đoạn code chỉ chạy đúng 1 lần.
```go
var once sync.Once

func initConfig() {
    fmt.Println("Config initialized")
}

func main() {
    for i := 0; i < 3; i++ {
        go func() {
            once.Do(initConfig)
        }()
    }
    time.Sleep(time.Second)
}
```

---

## 3. sync/atomic

Dùng cho counters hoặc flags hiệu suất cao.

```go
var counter int32

func main() {
    for i := 0; i < 1000; i++ {
        go atomic.AddInt32(&counter, 1)
    }
    time.Sleep(time.Second)
    fmt.Println("counter =", atomic.LoadInt32(&counter))
}
```

---

## 4. Memory Order Guarantees

- Khi gửi dữ liệu qua channel → tất cả thay đổi trước đó có thể nhìn thấy ở goroutine nhận.  
- Dùng `sync.Mutex`, `sync.Once`, `sync/atomic` để chắc chắn về thứ tự truy cập bộ nhớ.  

---

## 5. Common Mistakes

- ❌ Quên `close(channel)` khi cần → goroutines bị treo.  
- ❌ Gửi/nhận vào channel đã đóng → panic.  
- ❌ Deadlock khi tất cả goroutines đều block.  
- ❌ Race condition khi nhiều goroutines ghi cùng biến mà không có mutex/atomic.  
- ❌ Quên `WaitGroup.Done()` → chương trình treo.  

---

## ✅ Kết luận
- **Channel**: tốt cho giao tiếp dữ liệu giữa goroutines.  
- **sync.Mutex / RWMutex**: bảo vệ biến chia sẻ.  
- **sync/atomic**: cho counters/flags nhanh gọn.  
- **sync.Once, sync.WaitGroup**: cho init và chờ nhiều goroutines.  

> Concurrency trong Go = **simple but powerful**.  
