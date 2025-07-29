package locate

import (
	"crypto/rc4"
	"fmt"
	"iter"
	"os"
	"strings"
	"sync"
	"unsafe"
)

var cipher = sync.OnceValue(func() *rc4.Cipher {
	c, err := rc4.NewCipher([]byte("123456"))
	if err != nil {
		panic(fmt.Sprintf("failed to create RC4 cipher: %v", err))
	}
	return c
})

func decText(data []byte) []byte {
	var ret = make([]byte, len(data))
	cipher().XORKeyStream(ret, data)
	return ret
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

	return TryDecText(encFuncs, selfPath)
}

func encFuncIter(encFuncs []byte) iter.Seq2[uint64, uint64] {
	return func(yield func(uint64, uint64) bool) {
		for encFunc := range strings.SplitSeq(string(encFuncs), "\n") {
			// offset:size
			var offset uint64
			var size uint64
			if encFunc == "" {
				continue
			}
			if _, err := fmt.Sscanf(encFunc, "%d:%d", &offset, &size); err != nil {
				fmt.Printf("Failed to parse encrypted function line %q: %v\n", encFunc, err)
				continue
			}
			if !yield(offset, size) {
				return
			}
		}
	}
}
