package locate

import (
	"crypto/rc4"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func CopyToLocation(location uintptr, data []byte) {
	f := RawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_WRITE)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
}

func pageStart(ptr uintptr) uintptr {
	return ptr &^ (uintptr(syscall.Getpagesize() - 1))
}

func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := RawMemoryAccess(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}

func RawMemoryAccess(p uintptr, length int) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(p)), length)
}

var Key = []byte("123456")

func DecFunc(selfPath string) error {
	encFuncs, err := os.ReadFile(selfPath + ".func")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No encrypted functions found in", selfPath)
			return nil
		}
		return err
	}

	if len(encFuncs) == 0 {
		fmt.Println("no encrypted functions found")
		return nil
	}

	mappings, err := LocatePlugin(selfPath)
	if err != nil {
		return err
	}
	for _, mapping := range mappings {
		fmt.Println("Found mapping:", mapping.String())
	}

	var cipher, _ = rc4.NewCipher(Key)

	_ = cipher

	for encFunc := range strings.SplitSeq(string(encFuncs), "\n") {
		// offset:size
		var offset uint64
		var size uint64
		if _, err := fmt.Sscanf(encFunc, "%d:%d", &offset, &size); err != nil {
			continue
		}

		for _, mapping := range mappings {
			if offset < mapping.Offset {
				continue
			}
			if offset >= mapping.Offset+mapping.Size() {
				continue
			}
			off := offset - mapping.Offset + mapping.Start
			if off+size > mapping.End {
				continue
			}
			fmt.Printf("Decrypting function at offset %#x, size %#x\n", off, size)

			ciphertext := RawMemoryAccess(uintptr(off), int(size))
			plaintext := make([]byte, size)
			cipher.XORKeyStream(plaintext, ciphertext)
			CopyToLocation(uintptr(off), plaintext)
		}
	}
	return nil
}
