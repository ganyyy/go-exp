package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	Addr = "localhost:6380"
	DB = 10
	InitClient()

	const (
		Key   = "MyStream3"
		Group = "group1"
	)

	AddToStream(Key, "key1", "val1", "key2", "val2")
	AddToStream(Key, "key1", "val1", "key2", "val2")
	AddToStream(Key, "key1", "val1", "key2", "val2")

	ret, err := RangeStream(Key)
	assert.Nil(t, err)
	t.Logf("Range: %+v", ret)

	ret, err = ReadFromStream(Key)
	assert.Nil(t, err)
	t.Logf("Read: %+v", ret)

	cnt, err := StreamAck(Key, Group, ret[0].ID)
	assert.Nil(t, err)
	t.Logf("Ack: %+v", cnt)

	ret, err = RangeStream(Key)
	assert.Nil(t, err)
	t.Logf("Range: %+v", ret)
}
