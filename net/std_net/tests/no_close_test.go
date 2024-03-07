package tests

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoCloseTcpServer(t *testing.T) {
	var listener, err = net.Listen("tcp", ":9999")
	assert.NoError(t, err)
	client, err := listener.Accept()
	assert.NoError(t, err)
	t.Logf("client:%v", client.RemoteAddr())
	select {}
}
