package locate

/*

#include "parse_darwin.h"

*/
import "C"
import (
	"fmt"
	"unsafe"
)

func TryDecText(encFuncs []byte, selfPath string) error {
	// Implementation for Darwin

	imageSlide := imageSlide(selfPath)
	for offset, size := range encFuncIter(encFuncs) {
		fmt.Println("Decoding function at offset", offset, "with size", size)
		addr := imageSlide + uintptr(offset)
		if err := patchText(addr, decText(RawMemoryAccess(addr, int(size)))); err != nil {
			panic(fmt.Sprintf("failed to patch text at %x: %v", addr, err))
		}
	}

	return nil
}

func imageSlide(selfPath string) uintptr {
	selfPathC := C.CString(selfPath)
	defer C.free(unsafe.Pointer(selfPathC))
	addr := C.get_image_slide(selfPathC)
	if addr == 0 {
		panic(fmt.Sprintf("failed to get image slide for %s", selfPath))
	}
	return uintptr(addr)
}

func patchText(addr uintptr, data []byte) error {
	var errMsg *C.char
	cdata := (*C.char)(unsafe.Pointer(&data[0]))
	ret := C.patch_text(C.uintptr_t(addr), C.size_t(len(data)), cdata, &errMsg)
	if ret != 0 {
		return fmt.Errorf("patch_text failed: %s", C.GoString(errMsg))
	}
	return nil
}
