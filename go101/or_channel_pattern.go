package go101

import (
	"fmt"
	"sync"
	"time"
)

var sig = func(name string, d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		time.Sleep(d)

		if name == "B" {
			panic("test error")
		}
		fmt.Println("worker", name, "done after", d)
	}()
	fmt.Println("worker", name, "return")
	return ch
}

// smae promise race chỉ quan tâm thằng nào trả về kết quả trc (không quan tâm succes hay fail)
func Or(chs ...<-chan struct{}) <-chan struct{} {
	done := make(chan struct{})
	var once sync.Once

	for _, v := range chs {
		go func(v <-chan struct{}) {
			<-v
			once.Do(func() {
				close(done)
			})
		}(v)
	}
	return done
}

func HandleOrPattern() {
	<-Or(
		sig("A", time.Second*3),
		sig("B", time.Second*2),
		sig("C", time.Second*4),
	)
}
