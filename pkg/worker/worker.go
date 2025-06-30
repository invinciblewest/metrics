package worker

import (
	"context"
	"runtime"
	"sync"
)

// Job представляет собой функцию, которая принимает контекст и возвращает ошибку.
type Job func(ctx context.Context) error

// Pool представляет собой пул горутин, которые выполняют задачи (Job).
type Pool struct {
	jobCh chan Job
	wg    sync.WaitGroup
}

// NewPool создает новый экземпляр Pool с каналом для задач и группой ожидания.
func NewPool() *Pool {
	return &Pool{
		jobCh: make(chan Job, runtime.NumCPU()),
		wg:    sync.WaitGroup{},
	}
}

// Start запускает пул горутин.
// Он принимает контекст для управления временем выполнения и ограничение по количеству горутин (rateLimit).
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

// Wait блокирует выполнение до тех пор, пока все горутины в пуле не завершат свою работу.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// AddJob добавляет новую задачу (Job) в пул для выполнения.
func (p *Pool) AddJob(job Job) {
	p.jobCh <- job
}

// Stop останавливает пул, закрывая канал задач.
func (p *Pool) Stop() {
	close(p.jobCh)
}
