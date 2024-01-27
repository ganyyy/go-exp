package anydata

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"unsafe"
)

type Codec interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

type JsonCodec[T any] struct{}

// Encode implements the json.Marshaler interface.
func (a *JsonCodec[T]) Encode() ([]byte, error) {
	var pt = (*T)(unsafe.Pointer(a))
	return json.Marshal(pt)
}

// Decode implements the json.Unmarshaler interface.
func (a *JsonCodec[T]) Decode(data []byte) error {
	var pt = (*T)(unsafe.Pointer(a))
	return json.Unmarshal(data, pt)
}

func (a *JsonCodec[T]) ShowFields() {
	var pt = (*T)(unsafe.Pointer(a))
	var rv = reflect.ValueOf(pt).Elem()
	println(uintptr(unsafe.Pointer(pt)))
	for i := 0; i < rv.NumField(); i++ {
		println(rv.Type().Field(i).Name, rv.Field(i).UnsafeAddr())
	}
}

type GobCodec[T any] struct{}

func (a *GobCodec[T]) Encode() ([]byte, error) {
	var pt = (*T)(unsafe.Pointer(a))
	var buffer = bytes.NewBuffer(nil)
	var encoder = gob.NewEncoder(buffer)
	err := encoder.Encode(pt)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// implements the Codec interface.
func (a *GobCodec[T]) Decode(data []byte) error {
	var pt = (*T)(unsafe.Pointer(a))
	var buffer = bytes.NewBuffer(data)
	var decoder = gob.NewDecoder(buffer)
	return decoder.Decode(pt)
}

type Player struct {
	GobCodec[Player]
	Name    string
	Address string
	Age     int
}
