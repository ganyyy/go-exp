package gopool

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

var (
	ErrFull  = errors.New("Pool Work channel full!")
	ErrClose = errors.New("Pool isClose!")
)

type worker struct {
	f Runner
}

func (w *worker) init() {
}

type Pool struct {
	wait      sync.WaitGroup
	num       int
	closeChan chan struct{}
	workChan  chan *worker
	closed    int32
}

type Runner func()

func NewPool(num int) *Pool {
	return &Pool{
		num:       num,
		closeChan: make(chan struct{}),
		workChan:  make(chan *worker, 1024),
	}
}

func (p *Pool) add(w worker) error {
	w.init()
	select {
	case p.workChan <- &w:
		return nil
	default:
		return ErrFull
	}
}

func (p *Pool) Start() {
	p.wait.Add(p.num)
	for i := 0; i < p.num; i++ {
		go func(i int) {
			defer p.wait.Done()
			log.Printf("[%v] Start", i)
			for {
				select {
				case w := <-p.workChan:
					w.f()
				case <-p.closeChan:
					log.Printf("[%v] Done", i)
					return
				}
			}
		}(i)
	}
}

func (p *Pool) Close() {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return
	}
	close(p.closeChan)
	p.wait.Wait()
}

func (p *Pool) isClose() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

func (p *Pool) Run(r Runner) error {
	if p.isClose() {
		return ErrFull
	}
	return p.add(worker{f: r})

}
