package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

func TestCodec(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	var encoder = codec.NewEncoder(buf, &codec.MsgpackHandle{})
	var decoder = codec.NewDecoder(buf, &codec.MsgpackHandle{})

	type Value struct {
		_struct struct{} `codec:",omitempty"`

		A int
		B int
		C bool
	}

	var v1 Value
	//v1.A = 100
	var v2 Value

	_ = encoder.Encode(v1)
	t.Logf("%v", buf.String())

	_ = decoder.Decode(&v2)

	assert.Equal(t, v1, v2)

}
