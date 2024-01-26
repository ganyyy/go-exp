package main

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func TestNetConnect(t *testing.T) {
	// 系统套接字创建时会设置为非阻塞的, 这时候就可以通过context控制连接超时
	_, err := net.DialTimeout("tcp", "localhost:6666", time.Second)
	t.Logf("conn error:%v", err)

}

func TestNetMSSAndRecvBuffer(t *testing.T) {

	showFDInfo := func(from string, fd int) {
		recvBuffer, err := syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF)
		if err != nil {
			t.Errorf("%s getsockopt SO_RCVBUF error:%v", from, err)
			return
		}
		sendBuffer, err := syscall.GetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_SNDBUF)
		if err != nil {
			t.Errorf("%s getsockopt SO_SNDBUF error:%v", from, err)
			return
		}
		mss, err := syscall.GetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_MAXSEG)
		if err != nil {
			t.Errorf("%s getsockopt TCP_MAXSEG error:%v", from, err)
			return
		}
		t.Logf("%s fd:%d, recvBuffer:%d, sendBuffer:%d mss:%d", from, fd, recvBuffer, sendBuffer, mss)
	}

	var dialCfg = &net.Dialer{
		Control: func(network, address string, c syscall.RawConn) error {
			var fd uintptr
			c.Control(func(f uintptr) { fd = f })
			showFDInfo("before", int(fd))
			t.Logf("network:%s, address:%s, fd:%d", network, address, fd)
			return nil
		},
	}
	conn, err := dialCfg.Dial("tcp", "localhost:6666")
	if err != nil {
		t.Errorf("dial error:%v", err)
		return
	}
	row, err := conn.(*net.TCPConn).SyscallConn()
	if err != nil {
		t.Errorf("get syscall conn error:%v", err)
		return
	}
	err = row.Control(func(fd uintptr) { showFDInfo("after", int(fd)) })
	if err != nil {
		t.Errorf("control syscall conn error:%v", err)
		return
	}
}
