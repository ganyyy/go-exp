package main

import (
	"bytes"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

func TestCodec(t *testing.T) {
	type Value struct {
		Name  string `codec:"name"`
		Empty int    `codec:"-"`
	}

	var buf = bytes.NewBuffer(nil)
	var encoder = codec.NewEncoder(buf, &codec.MsgpackHandle{})
	var decoder = codec.NewDecoder(buf, &codec.MsgpackHandle{})

	assert.Nil(t, encoder.Encode(&Value{
		Name:  "1312312",
		Empty: 100,
	}))
	var v Value
	assert.Nil(t, decoder.Decode(&v))

	t.Logf("%+v", v)
}

type MyLocker struct {
	lock sync.Mutex
}

func (lock *MyLocker) Locker() func() {
	lock.lock.Lock()
	return lock.lock.Unlock
}

func TestLocker(t *testing.T) {
	var m MyLocker
	var cnt int64
	const NUM = 1000000
	var wait sync.WaitGroup

	var locker = func() func() {
		m.lock.Lock()
		return m.lock.Unlock
	}

	wait.Add(NUM)
	for i := 0; i < NUM; i++ {
		go func() {
			defer wait.Done()
			// m.lock.Lock()
			// defer m.lock.Unlock()
			defer locker()()
			cnt++
		}()
	}

	wait.Wait()

	assert.Equal(t, int64(NUM), atomic.LoadInt64(&cnt))
}
