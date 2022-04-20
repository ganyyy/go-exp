package msg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	t.Run("msg pack", func(t *testing.T) {
		var data = []byte("hello world!")
		var msg = NewMsg(MsgLogin, data)
		var pack = msg.Pack()

		var msg2, err = ReadMsg(pack)
		assert.Nil(t, err)
		assert.Equal(t, msg, msg2)
	})
}
