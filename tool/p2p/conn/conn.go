package conn

import (
	"p2p/msg"
	"syscall"
)

type Conn struct {
	fd int
}

func NewConn() (*Conn, error) {
	var sockFd int
	var err error
	sockFd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	return &Conn{fd: sockFd}, err
}

func (c *Conn) Fd() int {
	return c.fd
}

func (c *Conn) SendMsg(msg msg.Msg, to Addr) error {
	return SendMsg(c.Fd(), msg, to)
}

func (c *Conn) ReadMsg() (msg msg.Msg, addr Addr, err error) {
	return RecvMsg(c.Fd())
}

func SendMsg(fd int, msg msg.Msg, to Addr) error {
	return syscall.Sendmsg(fd, msg.Pack(), nil, to.ToSysInet4Addr(), 0)
}

func RecvMsg(fd int) (msg msg.Msg, addr Addr, err error) {
	var buf [1024]byte
	n, sockAddr, err := syscall.Recvfrom(fd, buf[:], 0)
	if err != nil {
		return
	}
	inet4Addr, _ := sockAddr.(*syscall.SockaddrInet4)
	addr.IP = inet4Addr.Addr
	addr.Port = uint16(inet4Addr.Port)
	err = msg.Unpack(buf[:n])
	return
}
