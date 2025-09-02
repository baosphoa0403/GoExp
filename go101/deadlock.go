package go101

import (
	"fmt"
	"sync"
	"time"
)

func DeadLockWg() {
	var m1, m2 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func() { // G1: lock m1 → m2
		defer wg.Done()
		m1.Lock()
		defer m1.Unlock()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("G1: try lock m2")
		m2.Lock()
		defer m2.Unlock()
		fmt.Println("G1 done")
	}()

	go func() { // G2: lock m2 → m1
		defer wg.Done()
		m2.Lock()
		defer m2.Unlock()
		time.Sleep(100 * time.Millisecond)
		fmt.Println("G2: try lock m1")
		m1.Lock()
		defer m1.Unlock()
		fmt.Println("G2 done")
	}()

	wg.Wait() // treo mãi (fatal: all goroutines are asleep - deadlock!)
}
