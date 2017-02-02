package util

import (
	"fmt"
)

type _taskch chan func()
type _poolch chan _taskch

type WorkerPool struct {
	pool    _poolch
	workers []*_worker
	started bool
}

type _worker struct {
	task _taskch
	pool _poolch
	quit chan struct{}
}

func (this *WorkerPool) Start(maxWorkers int) {
	this.started = true
	this.pool = make(_poolch, maxWorkers)
	this.workers = make([]*_worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		worker := new(_worker)
		worker.pool = this.pool
		worker.task = make(chan func(), 1)
		worker.quit = make(chan struct{}, 1)
		this.workers[i] = worker
		go worker.waitForJob()
	}
}

func (this *WorkerPool) Quit() {
	for _, worker := range this.workers {
		worker.quit <- struct{}{}
	}
}

func (this *WorkerPool) AddWork(task func()) (err error) {
	if !this.started {
		return fmt.Errorf("not started.")
	}
	go func() {
		taskch := <-this.pool
		taskch <- task
	}()

	return nil
}

func (this *_worker) waitForJob() {
	for {
		this.pool <- this.task

		select {
		case job := <-this.task:
			job()
		case <-this.quit:
			return
		}
	}

}
