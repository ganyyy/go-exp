package locate

/*
#cgo LDFLAGS: -ldl
#ifdef __linux__
#define _GNU_SOURCE
#endif
#include <dlfcn.h>
#include <stdlib.h>

void ____TUZeN7l3WTDytJk() {};

// C 包装函数，调用 dladdr 获取包含 GoAnchor 的库路径
const char* find_e1T1aA2Y_plugin_self_path() {
	Dl_info info;
	if (dladdr((void*)____TUZeN7l3WTDytJk, &info) == 0 || info.dli_fname == NULL) {
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
	selfPath := C.find_e1T1aA2Y_plugin_self_path()
	if selfPath == nil {
		return "", os.ErrNotExist
	}
	return C.GoString(selfPath), nil
}

func Locate() (string, error) {
	return locateSelf()
}
