package core

import (
	"sync"
	"time"
)

type GetTask func() Runner

type result struct {
	duration time.Duration
	status   int
	response interface{}
	err      error
}

type Bee struct {
	QPS           int
	Tasks         chan Runner
	Medium        chan Runner
	Duration      time.Duration
	Timeout       time.Duration
	ProducerCount int
	ConsumerCount int
	TaskHandler   GetTask
	results       chan *result
	allDone       chan bool
}

func (b Bee) Launcher() {
	defer close(b.Medium)
	defer close(b.Tasks)
	b.results = make(chan *result)
	b.allDone = make(chan bool)

	var wg sync.WaitGroup
	var rwg sync.WaitGroup

	go b.Scheduler(b.Duration)

	rwg.Add(1)
	go collectResults(b.results, b.allDone, &rwg)

	for i := 0; i < b.ConsumerCount; i++ {
		wg.Add(1)
		go b.worker(b.Timeout, &wg)
	}

	for i := 0; i < b.ProducerCount; i++ {
		wg.Add(1)
		go b.queen(b.Timeout, &wg)
	}

	wg.Wait()
	b.allDone <- true
	rwg.Wait()
}

func (b Bee) Scheduler(duration time.Duration) {
	interval := time.Duration(1e6/b.QPS) * time.Microsecond
	tick := time.Tick(interval)
	timer := time.NewTimer(duration)

	for {
		select {
		case <-tick:
			b.Medium <- b.TaskHandler()
		case <-timer.C:
			return
		}
	}
}

func (b Bee) run(task Runner, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	start := time.Now()
	waitGroup.Add(1)
	response, err := task.Run()
	duration := time.Now().Sub(start)
	b.results <- &result{duration: duration, status: 1, response: response, err: err}
}

func (b Bee) queen(timeout time.Duration, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for {
		select {
		case task := <-b.Medium:
			b.Tasks <- task
		case <-time.After(timeout):
			return
		}
	}
}

func (b Bee) worker(timeout time.Duration, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for {
		select {
		case task := <-b.Tasks:
			go b.run(task, waitGroup)
		case <-time.After(timeout):
			return
		}
	}
}
