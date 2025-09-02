# Go Optimization Guide 🚀

Tài liệu tóm tắt các cách **tối ưu hiệu năng trong Golang**, kèm ví dụ
minh họa.\
👉 Rule chung: **Đo trước, tối ưu sau** (`pprof`, `bench`).

------------------------------------------------------------------------

## 📌 Checklist

### 1. Algorithm & Data Structures

-   Chọn thuật toán & cấu trúc dữ liệu đúng.
-   Tránh `O(n²)` nếu có thể dùng `O(n log n)` hoặc `O(1)`.

**Ví dụ:**

``` go
// ❌ Tìm trong slice: O(n)
func containsBad(arr []int, x int) bool {
    for _, v := range arr {
        if v == x {
            return true
        }
    }
    return false
}

// ✅ Dùng map: O(1)
func containsGood(m map[int]struct{}, x int) bool {
    _, ok := m[x]
    return ok
}
```

------------------------------------------------------------------------

### 2. Memory & Allocation

-   Tránh tạo object/slice/string tạm nhiều lần.
-   Pre-allocate slice bằng `make`.
-   Dùng `strings.Builder` hoặc `sync.Pool` để giảm GC.

**Ví dụ:**

``` go
// ❌ Nối string tốn allocation
func concatBad(data []string) string {
    out := ""
    for _, s := range data {
        out += s
    }
    return out
}

// ✅ Dùng strings.Builder
func concatGood(data []string) string {
    var b strings.Builder
    b.Grow(len(data) * 10)
    for _, s := range data {
        b.WriteString(s)
    }
    return b.String()
}
```

------------------------------------------------------------------------

### 3. Concurrency Control

-   Không spawn goroutine vô hạn → dùng **worker pool**.
-   Channel nên có buffer hợp lý.
-   Luôn dùng `context` để cancel → tránh leak.

**Ví dụ:**

``` go
// ❌ Goroutine leak
func leak() chan int {
    ch := make(chan int)
    go func() {
        for { ch <- 1 }
    }()
    return ch
}

// ✅ Safe goroutine với context
func safe(ctx context.Context) chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for {
            select {
            case <-ctx.Done():
                return
            case ch <- 1:
            }
        }
    }()
    return ch
}
```

------------------------------------------------------------------------

### 4. I/O Optimization

-   Gom nhiều write nhỏ thành 1 write lớn.
-   Dùng `bufio.Writer` để giảm syscall.

**Ví dụ:**

``` go
// ❌ Nhiều syscall
for i := 0; i < 1000; i++ {
    f.Write([]byte("hello\n"))
}

// ✅ Gom buffer
w := bufio.NewWriter(f)
for i := 0; i < 1000; i++ {
    w.WriteString("hello\n")
}
w.Flush()
```

------------------------------------------------------------------------

### 5. Profiling-driven

Dùng benchmark & pprof để tìm bottleneck:

``` bash
go test -bench . -benchmem
go test -cpuprofile cpu.out -memprofile mem.out
go tool pprof cpu.out
```

------------------------------------------------------------------------

### 6. Avoid reflection & interface overhead

-   Dùng **generics** thay cho `interface{}` + `reflect`.
-   Hạn chế marshal/unmarshal bằng reflect.

**Ví dụ:**

``` go
// ❌ Reflect chậm
val := reflect.ValueOf(myStruct)
fmt.Println(val.FieldByName("Name"))

// ✅ Access trực tiếp
fmt.Println(myStruct.Name)
```

------------------------------------------------------------------------

### 7. Trade-off

-   Đừng tối ưu quá sớm.
-   Code dễ đọc, maintain dễ thường quan trọng hơn micro-optimize.
-   Nguyên tắc: **Correct → Clean → Measurable → Optimize**.

------------------------------------------------------------------------

## ✅ Kết luận

-   Tập trung **algorithm + memory allocation** trước tiên.
-   Dùng `pprof` để optimize đúng chỗ.
-   Luôn thiết kế hệ thống để **tránh leak goroutines và I/O syscall
    thừa**.

Happy hacking with Go 🚀
