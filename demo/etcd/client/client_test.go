package client

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func TestSetGet(t *testing.T) {
	Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:2379"},
		Root:      "/root",
	})

	assert.Nil(t, initError)

	Put("/test/v1", "123")
	Put("/test/v2", "456")
	Put("/test/v3", "789")

	var ret, _ = Get("/test", true)

	for _, kv := range ret {
		t.Logf("%+v", kv)
	}

	const WatchKey = "/test/temporary"

	var ctx, cancel = context.WithCancel(context.Background())
	var ctx2, cancel2 = context.WithCancel(context.Background())

	Watch(ctx, WatchKey, false)
	Watch(ctx2, "/test", true)

	for i := range 3 {
		Put(WatchKey, "Val"+strconv.Itoa(i))
		Put(WatchKey+"/last", "Val"+strconv.Itoa(i))
		time.Sleep(time.Second)
	}
	cancel()
	time.Sleep(time.Second)
	cancel2()
	time.Sleep(time.Second)
}

func Range(n int) []struct{} {
	return make([]struct{}, n)
}

func TestCAS(t *testing.T) {
	Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:2379"},
		Root:      "/root",
	})

	defaultClient.Delete(context.Background(), "/test/cas")

	valid, newVal, err := CAS("/test/cas", "456", "")
	assert.Nil(t, err)
	assert.True(t, valid)
	assert.Equal(t, "456", newVal)

	valid, newVal, err = CAS("/test/cas", "789", "456")
	assert.Nil(t, err)
	assert.True(t, valid)
	assert.Equal(t, "789", newVal)

	valid, newVal, err = CAS("/test/cas", "459", "788")
	assert.Nil(t, err)
	assert.False(t, valid)
	assert.Equal(t, "789", newVal)
}

func TestNX(t *testing.T) {
	Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:2379"},
		Root:      "/root",
	})

	defaultClient.Delete(context.Background(), "/test/nx")

	valid, err := SetNX("/test/nx", "456")
	assert.Nil(t, err)
	assert.True(t, valid)

	valid, err = SetNX("/test/nx", "456")
	assert.Nil(t, err)
	assert.False(t, valid)
}

func TestDistributedLock(t *testing.T) {

	const Lock = "/test/lock"

	Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:2379"},
		Root:      "/root",
	})
	assert.Nil(t, initError)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		unlock, err := DistributedLock(ctx, Lock)
		if !assert.Nil(t, err) {
			t.Logf("1 lock failed")
			return
		}
		defer unlock()
		time.Sleep(8 * time.Second)
		t.Logf("1 lock success")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		unlock, err := DistributedLock(ctx, Lock)
		if !assert.Nil(t, err) {
			t.Logf("2 lock failed")
			return
		}
		defer unlock()
		t.Logf("2 lock success")
		time.Sleep(time.Second * 5)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		unlock, err := DistributedLock(ctx, Lock)
		if !assert.Nil(t, err) {
			t.Logf("3 lock failed")
			return
		}
		defer unlock()
		t.Logf("3 lock success")
	}()

	wg.Wait()
}

func TestSessionLock(t *testing.T) {
	require.NoError(t, Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:2379"},
		Root:      "/root",
	}))

	session, err := concurrency.NewSession(defaultClient.Client, concurrency.WithTTL(10))
	require.Nil(t, err)
	defer session.Close()

	const Lock = "/test/session/lock"

	unlock, err := DistributedLockWithSession(context.Background(), session, Lock)
	require.Nil(t, err)
	defer unlock()
	t.Logf("lock success")
	<-time.After(5 * time.Second)
	t.Logf("lock exit")
}
