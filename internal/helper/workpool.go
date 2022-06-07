package helper

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

func NewGoroutinePool(workerSize int, poolSize int) *GoroutinePool {
	gp := &GoroutinePool{
		queue: make(chan work, poolSize),
	}

	gp.AddWorkers(workerSize)
	return gp
}

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

func (gp *GoroutinePool) Close() {
	close(gp.queue)
	gp.wg.Wait()
}

func (gp *GoroutinePool) ScheduleWork(fn WorkFunc) {
	gp.queue <- work{fn}
}
