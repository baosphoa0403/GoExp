package utils

import (
	"fmt"
	"sync"
)

type Pool struct {
	workerCount int
	jobs        chan func() // channel function because not take care input output
	wg          sync.WaitGroup
}

func NewPool(workerCount int) *Pool {
	pool := Pool{
		workerCount: workerCount,
		wg:          sync.WaitGroup{},
		jobs:        make(chan func(), 10),
	}

	pool.start()

	return &pool
}

func (p *Pool) start() {
	for index := 0; index < p.workerCount; index++ {
		go func(workerId int) {
			for job := range p.jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("panic recovered", r)
						}
					}()
					job()
				}()
			}
		}(index)
	}
}

func (p *Pool) Submit(job func()) {
	p.wg.Add(1)
	p.jobs <- func() {
		defer p.wg.Done()
		job()
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Stop() {
	close(p.jobs)
}
