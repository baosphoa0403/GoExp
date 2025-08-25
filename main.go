package main

import "giabao.com/go101"

func main() {
	// goroutine.Goroutine()
	// pdf()
	// lock.MutexLock()

	// pool := utils.NewPool(10)

	// for i := 0; i < 15; i++ {
	// 	jobID := i
	// 	pool.Submit(func() {
	// 		fmt.Printf("Job %d running\n", jobID)
	// 		time.Sleep(500 * time.Millisecond)

	// 		if jobID == 7 {
	// 			panic("Job 7 bị lỗi nè")
	// 		}

	// 		fmt.Printf("Job %d done\n", jobID)
	// 	})
	// }
	// pool.Wait()
	// pool.Stop()

	// concurrency.MainBoring()
	// concurrency.MainChannelBoring()
	// concurrency.TestCtx()
	// concurrency.TestErrorGroup()

	// go 101
	// go101.FanOutGoroutineProducer()
	// go101.FanOutMainProducer()
	go101.FanIn()

}
