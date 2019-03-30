package main

import (
	"errors"
	"sync"
)

type Dispatcher struct {
	name           string
	workers        []*worker
	freeWorkerPool chan *worker
	taskQueue      chan interface{}
	taskFunc       TaskFunction
	stopChan       chan struct{}
	wg             sync.WaitGroup
	mode           int
}

type TaskFunction func(interface{}) error

const (
	QueueMode int = iota
	ContinualMode
)

func NewDispatcher(name string, taskFunc TaskFunction, workerNum int, mode int) *Dispatcher {
	dispatcher := &Dispatcher{
		name:           name,
		freeWorkerPool: make(chan *worker, workerNum),
		taskQueue:      make(chan interface{}),
		stopChan:       make(chan struct{}),
		taskFunc:       taskFunc,
		mode:           mode,
	}
	dispatcher.workers = make([]*worker, workerNum)
	for i := 0; i < workerNum; i++ {
		w := worker{
			id:         i,
			dispatcher: dispatcher,
			jobPool:    make(chan interface{}),
		}
		dispatcher.workers[i] = &w
	}
	return dispatcher
}

func (dispatcher *Dispatcher) Start() {
	for _, worker := range dispatcher.workers {
		dispatcher.freeWorkerPool <- worker
	}

	go func() {
		for {
			worker := <-dispatcher.freeWorkerPool
			job, ok := <-dispatcher.taskQueue
			if !ok {
				return
			}
			logger.WithWorkerInfo(worker).Debugln("job fetched")
			go worker.work(&job)
		}
	}()

	if dispatcher.mode == ContinualMode {
		go func() {
			for {
				if err := dispatcher.Add(struct{}{}); err != nil {
					return
				}
			}
		}()
	}
}

func (dispatcher *Dispatcher) Add(job interface{}) error {
	dispatcher.wg.Add(1)
	for {
		select {
		case <-dispatcher.stopChan:
			dispatcher.wg.Done()
			return errors.New("job queue not accepted")
		default:
			select {
			case dispatcher.taskQueue <- job:
				return nil
			default:
			}
		}
	}
}

func (dispatcher *Dispatcher) Stop() {
	close(dispatcher.stopChan) // taskQueueに送信させない
	dispatcher.wg.Wait()       // 全てのtaskが終わるまで待つ
	close(dispatcher.taskQueue)
}

type worker struct {
	id         int
	dispatcher *Dispatcher
	jobPool    chan interface{}
}

func (w *worker) work(job *interface{}) {
	defer w.dispatcher.wg.Done()

	err := w.dispatcher.taskFunc(job)
	w.dispatcher.freeWorkerPool <- w
	if err != nil {
		logger.WithWorkerInfo(w).Errorln(err)
		return
	}
	logger.WithWorkerInfo(w).Debugln("job finished")
}
