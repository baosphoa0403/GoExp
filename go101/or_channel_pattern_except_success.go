package go101

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	data  any
	error error
}

type MyData struct {
	Name string
}

var sigtask = func(name string, d time.Duration) <-chan Result {
	ch := make(chan Result)
	go func() {
		defer func() {
			if v := recover(); v != nil {
				ch <- Result{
					data:  struct{}{},
					error: fmt.Errorf("error panic", v),
				}
			}
			close(ch)
		}()

		time.Sleep(d)
		if name == "B" {
			panic("test error")
		}
		ch <- Result{
			data: MyData{
				Name: "gia bao",
			},
			error: nil,
		}

		fmt.Println("worker", name, "done after", d)
	}()
	return ch
}

// smae promise any chỉ quan tâm thằng nào trả về kết quả success
func OrSuccess(chs ...<-chan Result) <-chan Result {
	done := make(chan Result)
	var once sync.Once

	for _, v := range chs {
		go func(v <-chan Result) {
			data := <-v

			if data.error == nil {
				once.Do(func() {
					done <- Result{
						data:  data.data,
						error: nil,
					}
					close(done)
				})
			}

		}(v)
	}
	return done
}

func HandleOrPatternAny() {
	data := <-OrSuccess(
		sigtask("A", time.Second*1),
		sigtask("B", time.Second*2),
		sigtask("C", time.Second*4),
	)
	fmt.Println("data: ", data)
}
