package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	Init(&EtcdConfig{
		Endpoints: []string{"http://localhost:2379"},
		Root:      "/root",
	})

	assert.Nil(t, initError)

	defer Stop()

	Put("/test/v1", "123")
	Put("/test/v2", "456")
	Put("/test/v3", "789")

	var ret, _ = Get("/test", true)

	assert.Equal(t, len(ret), 3)

	for _, kv := range ret {
		t.Logf("%+v", kv)
	}

}
