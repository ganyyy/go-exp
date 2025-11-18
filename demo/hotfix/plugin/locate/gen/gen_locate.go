package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const Source = `package locate

/*
#cgo LDFLAGS: -ldl
#ifdef __linux__
#define _GNU_SOURCE
#endif
#include <dlfcn.h>
#include <stdlib.h>

void f1_%s() {};

// C 包装函数，调用 dladdr 获取包含 GoAnchor 的库路径
const char* f2_%s() {
	Dl_info info;
	if (dladdr((void*)f1_%s, &info) == 0 || info.dli_fname == NULL) {
		return NULL;
	}
	return info.dli_fname;
}
*/
import "C"
import (
	"os"
)

func locateSelf() (string, error) {
	selfPath := C.f2_%s()
	if selfPath == nil {
		return "", os.ErrNotExist
	}
	return C.GoString(selfPath), nil
}

func Locate() (string, error) {
	return locateSelf()
}
`

func main() {
	var buf1 [16]byte
	var buf2 [16]byte

	rand.Read(buf1[:])
	rand.Read(buf2[:])

	fn1 := hex.EncodeToString(buf1[:])
	fn2 := hex.EncodeToString(buf2[:])

	source := fmt.Sprintf(Source, fn1, fn2, fn1, fn2)

	fmt.Print(source)
}
