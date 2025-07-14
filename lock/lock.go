package lock

import (
	"fmt"
	"sync"
)

var (
	wg             sync.WaitGroup
	counter        int
	counterNoMutex int
	mutex          sync.Mutex
)

func increment() {
	defer wg.Done()
	mutex.Lock()
	counter++
	mutex.Unlock()
}

func incrementNoLock() {
	defer wg.Done()
	counterNoMutex++
}

func MutexLock() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go increment()

		// wg.Add(1)
		// go incrementNoLock()
	}

	wg.Wait()
	fmt.Println("final counter: ", counter)
	fmt.Println("final counterNoMutex: ", counterNoMutex)
}
