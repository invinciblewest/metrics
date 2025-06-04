package worker

import (
	"context"
	"runtime"
	"sync"
)

type Job func(ctx context.Context) error

type Pool struct {
	jobCh chan Job
	wg    sync.WaitGroup
}

func NewPool() *Pool {
	return &Pool{
		jobCh: make(chan Job, runtime.NumCPU()),
		wg:    sync.WaitGroup{},
	}
}

func (p *Pool) Start(ctx context.Context, rateLimit int) chan error {
	errorCh := make(chan error)

	for i := 0; i < rateLimit; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobCh:
					if !ok {
						return
					}
					if err := job(ctx); err != nil {
						errorCh <- err
					}
				}
			}
		}()
	}

	go func() {
		p.wg.Wait()
		close(errorCh)
	}()

	return errorCh
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) AddJob(job Job) {
	p.jobCh <- job
}

func (p *Pool) Stop() {
	close(p.jobCh)
}
