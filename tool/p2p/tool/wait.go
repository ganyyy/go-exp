package tool

import "sync"

type WaitGroup struct {
	group sync.WaitGroup
}

func (wait *WaitGroup) Do(f func()) {
	wait.group.Add(1)
	go func() {
		defer wait.group.Done()
		f()
	}()
}

func (wait *WaitGroup) Wait() {
	wait.group.Wait()
}
