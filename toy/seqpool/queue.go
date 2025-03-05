package seqpool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Queue[T any] struct {
	queue  []T // TODO 替换成 ring buffer or list
	lock   sync.Mutex
	signal chan struct{}
	n      chan struct{}
}

func NewQueue[T any]() *Queue[T] {
	n := make(chan struct{})
	close(n)
	return &Queue[T]{
		queue:  make([]T, 0),
		signal: make(chan struct{}, 1),
		n:      n,
	}
}

func (q *Queue[T]) Signal() <-chan struct{} {
	select {
	case <-q.signal:
		return q.n
	default:
	}
	if q.Len() > 0 {
		return q.n
	}
	return q.signal
}

func (q *Queue[T]) Push(v T) {
	q.lock.Lock()
	q.queue = append(q.queue, v)
	q.lock.Unlock()
	select {
	case q.signal <- struct{}{}:
	default:
	}
}

func (q *Queue[T]) Pop() (e T, ok bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.queue) == 0 {
		return
	}
	e = q.queue[0]
	var empty T
	q.queue[0] = empty
	q.queue = q.queue[1:]
	return e, true
}

func (q *Queue[T]) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.queue)
}

func (q *Queue[T]) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue = nil
}

type InfiniteQueue[T any] struct {
	buffer   []T
	inChan   chan T
	outChan  chan T
	stopChan chan struct{}
	isStop   atomic.Bool
}

func NewInfiniteQueue[T any]() *InfiniteQueue[T] {
	in := make(chan T, 1)
	out := make(chan T, 1)

	queue := &InfiniteQueue[T]{
		buffer:  make([]T, 0),
		inChan:  in,
		outChan: out,
	}

	go func() {
	end:
		for {
			var item T
			var out chan T
			var in = queue.inChan
			if len(queue.buffer) > 0 {
				val := queue.buffer[0]
				queue.buffer[0] = item
				item = val
				queue.buffer = queue.buffer[1:]
				out = queue.outChan
			}
			select {
			case <-queue.stopChan:
				break end
			case out <- item:
			case v := <-in:
				select {
				case <-queue.stopChan:
					break end
				default:
				}
				queue.buffer = append(queue.buffer, v)
			}
		}
	end2:
		for {
			select {
			case v := <-queue.inChan:
				queue.buffer = append(queue.buffer, v)
			default:
				break end2
			}
		}
		for _, v := range queue.buffer {
			queue.outChan <- v
		}
		close(queue.outChan)
	}()

	return queue
}

func (q *InfiniteQueue[T]) Enqueue(ctx context.Context, v T) bool {
	if q.isStop.Load() {
		return false
	}
	select {
	case <-ctx.Done():
		return false
	case <-q.stopChan:
		return false
	case q.inChan <- v:
		return !q.isStop.Load()
	}
}

func (q *InfiniteQueue[T]) Out() <-chan T {
	return q.outChan
}

func (q *InfiniteQueue[T]) Stop() {
	if q.isStop.CompareAndSwap(false, true) {
		close(q.stopChan)
	}
}
