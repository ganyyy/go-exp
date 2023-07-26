package mutex

import (
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestMutex(t *testing.T) {
	type Mutex struct {
		state int32
		sema  uint32
	}

	var mutex sync.Mutex
	var sm = (*Mutex)(unsafe.Pointer(&mutex))
	mutex.Lock()
	go func() {
		var ticker = time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				t.Logf("sm:%+v", sm)
			}
		}
		// mutex.Unlock()
	}()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		go func() {
			mutex.Lock()
		}()
	}

}

func Test_SyncConv(t *testing.T) {

	var lock sync.Mutex
	var cond = sync.NewCond(&lock)

	var cnt int

	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				lock.Lock()
				for cnt == 0 {
					cond.Wait()
				}
				t.Log(i, "wait cnt", cnt)
				cnt--
				lock.Unlock()
			}
		}(i)
	}

	go func() {
		for {
			lock.Lock()
			cond.Broadcast()
			cnt++
			t.Log("signal")
			lock.Unlock()
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 10)
}
