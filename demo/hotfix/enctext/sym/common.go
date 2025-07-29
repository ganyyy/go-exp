package sym

import (
	"crypto/rc4"
	"fmt"
	"strconv"
)

type EncFunc func(offset, size uint64) error

func GenRC4Enc(data, key []byte) EncFunc {
	var cipher, err = rc4.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("failed to create RC4 cipher: %v", err))
	}
	return func(offset, size uint64) error {
		if offset+size > uint64(len(data)) {
			return fmt.Errorf("offset %d and size %d exceed data length %d", offset, size, len(data))
		}
		cipher.XORKeyStream(data[offset:offset+size], data[offset:offset+size])
		return nil
	}
}

func AppendEncFunc(encFuncs []byte, offset, size uint64) []byte {
	var _buf [64]byte
	var buf = _buf[:0]
	buf = strconv.AppendUint(buf, offset, 10)
	buf = append(buf, ':')
	buf = strconv.AppendUint(buf, size, 10)
	buf = append(buf, '\n')
	return append(encFuncs, buf[:]...)
}
