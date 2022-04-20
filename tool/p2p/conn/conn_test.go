package conn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddr(t *testing.T) {
	t.Run("AddrParse", func(t *testing.T) {
		var buf = []byte("192.169.0.1:9980")
		var addr Addr
		assert.Nil(t, addr.FromBytes(buf))

		var str = addr.String()
		t.Logf(str)

		var bs = addr.Bytes()
		var addr2 Addr
		addr2.FromBytes(bs)
		assert.Equal(t, addr, addr2)
	})

	t.Run("ParseAddr", func(t *testing.T) {
		var src = "127.0.0.1:9999"
		t.Log((&Addr{}).FromString(src))
	})
}
