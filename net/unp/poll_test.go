package unp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

func TestUnixPoll(t *testing.T) {
	var socket, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	assert.NoError(t, err)
	err = unix.SetNonblock(socket, true)
	assert.NoError(t, err)

	err = unix.Connect(socket, &unix.SockaddrInet4{Port: 9999, Addr: [4]byte{127, 0, 0, 1}})
	assert.Equal(t, err, unix.EINPROGRESS)

	var event unix.PollFd
	event.Fd = int32(socket)
	event.Events = unix.POLLOUT
	var events = make([]unix.PollFd, 1)
	events[0] = event
	n, err := unix.Poll(events, 10000)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, unix.POLLOUT, int(events[0].Revents))
	unix.Write(socket, []byte("hello\n"))
}
