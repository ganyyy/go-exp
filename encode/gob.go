package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func main() {

	var buf = bytes.NewBuffer(nil)
	var encoder = gob.NewEncoder(buf)
	var decoder = gob.NewDecoder(buf)

	var src interface{} = &struct {
		Valid int8
	}{}

	var vv = true

	var val interface{} = true

	_ = encoder.Encode(val)

	fmt.Println(buf.Bytes())

	//buf.Reset()
	//_ = encoder.Encode(src)
	//
	//fmt.Println(buf.Bytes())

	var err = decoder.Decode(src)
	fmt.Println(err)

	vv, ok := src.(bool)

	fmt.Println(vv, ok)
}
