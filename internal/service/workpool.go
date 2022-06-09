package service

import "sync"

type WorkFunc interface {
	Run()
}

type GoroutinePool struct {
	queue chan work
	wg    sync.WaitGroup
}

type work struct {
	fn WorkFunc
}

// NewGoroutinePool initialize the worker pool and adds the given workers
func NewGoroutinePool(workerSize int, poolSize int) *GoroutinePool {
	gp := &GoroutinePool{
		queue: make(chan work, poolSize),
	}

	gp.AddWorkers(workerSize)
	return gp
}

// AddWorkers add a new worker
func (gp *GoroutinePool) AddWorkers(size int) {
	gp.wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			for job := range gp.queue {
				job.fn.Run()
			}
			gp.wg.Done()
		}()
	}
}

// Close closes the queue channel
func (gp *GoroutinePool) Close() {
	close(gp.queue)
	gp.wg.Wait()
}

// ScheduleWork add a job to the queue of the wp
func (gp *GoroutinePool) ScheduleWork(fn WorkFunc) {
	gp.queue <- work{fn}
}
