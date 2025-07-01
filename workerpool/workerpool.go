package workerpool

import (
	"AppointmentSummary_Assignment_Code/model"
	"sync"
)

var (
	Wp        model.WorkerPool
	WaitGroup sync.WaitGroup
)

// create a new worker pool
func NewWorkerPool(poolSize int, bufferSize int) {
	Wp.MaxWorker = poolSize
	Wp.QueuedTaskC = make(chan func(), bufferSize)
}

func Run() {
	for i := 0; i < Wp.MaxWorker; i++ {
		go func() {
			for task := range Wp.QueuedTaskC {
				task()
				WaitGroup.Done()
			}
		}()
	}
}

func AddTask(task func()) {
	Wp.QueuedTaskC <- task
}

func CloseWorkerPool() {
	close(Wp.QueuedTaskC)
}
