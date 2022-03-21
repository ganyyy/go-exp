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
