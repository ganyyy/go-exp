package client

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:3379"},
		Root:      "/root",
	})

	assert.Nil(t, initError)

	defer Stop()

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

	for i := range Range(10) {
		Put(WatchKey, "Val"+strconv.Itoa(i))
		Put(WatchKey+"/last", "Val"+strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}

func Range(n int) []struct{} {
	return make([]struct{}, n)
}
