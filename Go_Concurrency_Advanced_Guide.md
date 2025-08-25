# 🚀 Go Concurrency Advanced Guide

Tài liệu này bổ sung các kiến thức **nâng cao** về goroutines và concurrency trong Go, sau khi bạn đã nắm cơ bản.

---

## 1. Context Package

### 📌 Ý nghĩa
- Dùng để **hủy goroutine** khi hết thời gian, hết request hoặc không cần thiết nữa.  
- Thường dùng trong **HTTP server, gRPC, worker pool**.

### 📌 Các hàm chính
- `context.Background()` → context gốc.  
- `context.WithCancel(parent)` → có thể gọi cancel để dừng.  
- `context.WithTimeout(parent, d)` → tự động hủy sau `d`.  
- `context.WithDeadline(parent, t)` → tự động hủy sau thời điểm `t`.  

### 📌 Ví dụ
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

ch := make(chan string)
go func() {
    time.Sleep(3 * time.Second)
    ch <- "done"
}()

select {
case <-ch:
    fmt.Println("Received result")
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

---

## 2. Select nâng cao

### 📌 Ý nghĩa
- Lắng nghe nhiều channel cùng lúc.  
- Tránh deadlock hoặc chờ quá lâu.

### 📌 Ví dụ multiplexing
```go
select {
case msg1 := <-ch1:
    fmt.Println("Got", msg1)
case msg2 := <-ch2:
    fmt.Println("Got", msg2)
default:
    fmt.Println("No message yet")
}
```

👉 `default` giúp non-blocking.  
👉 Có thể kết hợp với `context.Done()` để hủy.  

---

## 3. Concurrency Patterns Nâng Cao

### 🔹 Bounded Worker Pool
Giới hạn số goroutines chạy song song.
```go
jobs := make(chan int, 100)
results := make(chan int, 100)

for w := 1; w <= 3; w++ {
    go func(id int) {
        for j := range jobs {
            results <- j * 2
        }
    }(w)
}

for j := 1; j <= 10; j++ {
    jobs <- j
}
close(jobs)

for i := 1; i <= 10; i++ {
    fmt.Println(<-results)
}
```

### 🔹 Pipeline nâng cao (multi-stage)
```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    c := gen(2, 3, 4)
    out := sq(c)
    for n := range out {
        fmt.Println(n)
    }
}
```

### 🔹 Fan-out + Fan-in với timeout
- Nhiều goroutines xử lý dữ liệu.  
- Gom kết quả vào 1 channel.  
- Hủy bằng `context.WithTimeout`.  

---

## 4. Scheduler & Performance

### 📌 Mô hình M-P-G
- **M (Machine)** = OS thread.  
- **P (Processor)** = scheduler logic trong Go.  
- **G (Goroutine)** = lightweight thread.  

### 📌 GOMAXPROCS
- Mặc định bằng số CPU logic.  
- Có thể chỉnh:  
```go
runtime.GOMAXPROCS(4) // chỉ dùng 4 cores
```

### 📌 Performance tips
- Hạn chế tạo goroutine vô hạn (leak).  
- Dùng `sync.Pool` để tái sử dụng object.  
- Benchmark với `go test -bench .`.  

---

## 5. Debugging & Tools

### 🔎 Race Detector
```bash
go run -race main.go
```
→ Phát hiện race condition.  

### 🔎 Profiling
```bash
go test -bench . -benchmem
go tool pprof
```
→ Kiểm tra CPU, memory.  

### 🔎 Tracing
- `runtime/trace` để phân tích scheduler và goroutines.  
- Kết hợp với `pprof` để tối ưu hiệu suất.  

---

## ✅ Kết luận

- **Context**: quản lý vòng đời goroutines.  
- **Select**: nghe nhiều channel.  
- **Patterns nâng cao**: worker pool, pipeline nhiều tầng, fan-in/out có timeout.  
- **Scheduler & Performance**: hiểu M-P-G, điều chỉnh GOMAXPROCS.  
- **Debugging**: race detector, pprof, trace.  

👉 Với những kiến thức này, bạn đã đi từ **cơ bản → nâng cao → expert** trong Go concurrency.
