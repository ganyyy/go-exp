package client

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	Watch(WatchKey, false)
	Watch("/test", true)

	for i := range Range(5) {
		Put(WatchKey, "Val"+strconv.Itoa(i))
		Put(WatchKey+"/last", "Val"+strconv.Itoa(i))
		time.Sleep(time.Second)
	}
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
