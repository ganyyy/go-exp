package sync

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

type Data struct {
	Value int
}

func (d *Data) Close(i int) {
	log.Printf("%v close data %+v", i, d.Value)
}

var sm sync.Map

var success sync.Map
var fail sync.Map
var closed sync.Map

func SafeStoreAndDeleteOld(key int, value *Data) error {
	log.Printf("add key %v, value %+v", key, value)
	if old, loaded := sm.LoadAndDelete(key); loaded {
		closed.Store(old.(*Data).Value, true)
		old.(*Data).Close(value.Value)
	}
	if old, loaded := sm.LoadOrStore(key, value); loaded {
		fail.Store(old.(*Data).Value, true)
		return fmt.Errorf("key %v already exists, old %+v", key, old)
	}
	success.Store(value.Value, true)
	return nil
}

// TestSafeLoadAndDelete 测试并发调用SafeLoadAndDelete, 保证不会误删不同goroutine添加的key相同的数据
func TestSafeLoadAndDelete(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int, t *testing.T) {
			defer wg.Done()
			if err := SafeStoreAndDeleteOld(1, &Data{Value: i}); err != nil {
				log.Printf("index %v err %+v", i, err)
			}
		}(i, t)
	}
	wg.Wait()

	sm.Range(func(key, value interface{}) bool {
		log.Printf("key %v, value %+v", key, value)
		return true
	})

	success.Range(func(key, value interface{}) bool {
		log.Printf("success key %v, value %+v", key, value)
		return true
	})

	fail.Range(func(key, value interface{}) bool {
		log.Printf("fail key %v, value %+v", key, value)
		return true
	})

	closed.Range(func(key, value interface{}) bool {
		log.Printf("closed key %v, value %+v", key, value)
		return true
	})

	closed.Swap(1, false)
	closed.CompareAndSwap(1, true, false)
}
