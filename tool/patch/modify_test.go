package patch

import (
	"testing"
	"unsafe"
)

func TestModifyString(t *testing.T) {
	var s = "123456"

	addr := uintptr(unsafe.Pointer(unsafe.StringData(s)))
	t.Logf("string %q address %x", s, addr)

	copyToLocation(addr, []byte("143456"))
	t.Logf("string %q address %x", s, addr)
}
