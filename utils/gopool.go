package utils

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Pool struct {
	jobs        chan func() // channel function because not take care input output
	wg          sync.WaitGroup
	mu          sync.Mutex
	workerCount int
	closed      bool
}

func NewPool(workerCount int, queueSize int) *Pool {
	pool := Pool{
		workerCount: workerCount,
		wg:          sync.WaitGroup{},
		jobs:        make(chan func(), queueSize),
	}

	pool.start()

	return &pool
}

func (p *Pool) start() {
	for index := 0; index < p.workerCount; index++ {
		go p.worker(index)
	}
}

func (p *Pool) worker(index int) {
	for job := range p.jobs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[Worker %d] panic recovered: %v\n", index, r)
				}
			}()
			job()
		}()
	}
}

func (p *Pool) Submit(job func()) {
	p.wg.Add(1)
	p.jobs <- func() {
		defer p.wg.Done()
		job()
	}
}

func (p *Pool) SubmitWithContext(ctx context.Context, job func(ctx context.Context)) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return errors.New("pool is closed")
	}
	p.mu.Unlock()

	p.wg.Add(1)

	wrappedJob := func() {
		defer p.wg.Done()

		select {
		case <-ctx.Done():
			// cancel task
			fmt.Println("⚠ Job bị huỷ: ", ctx.Err())
		default:
			job(ctx)
		}
	}

	select {
	case p.jobs <- wrappedJob:
		return nil
	default:
		// Queue đầy → cancel job ngay
		p.wg.Done()
		return errors.New("job queue full")
	}

}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Stop() {

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()

	p.Wait()
	close(p.jobs)
}
