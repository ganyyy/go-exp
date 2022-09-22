package main

import (
	"syscall"
	"testing"
	"time"
)

func TestSyscall(t *testing.T) {
	t.Run("Socket", func(t *testing.T) {

		const PROTO = syscall.IPPROTO_IP

		sk, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, PROTO)
		t.Log(err)

		addr := &syscall.SockaddrInet4{
			Port: 8899,
			Addr: [4]byte{127, 0, 0, 1},
		}

		// syscall.SetsockoptInt(sk, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)

		syscall.Bind(sk, addr)

		err = syscall.Listen(sk, 100)
		if err != nil {
			t.Logf(err.Error())
		}

		handle := func(sock int) {
			var buf [128]byte
			n, err := syscall.Read(sock, buf[:])
			if err != nil {
				return
			}
			t.Logf("sock:%v data:%v", sock, string(buf[:n]))
		}

		go func() {
			time.Sleep(time.Second)
			cli, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP)

			err := syscall.Connect(cli, addr)
			if err != nil {
				t.Logf(err.Error())
				return
			}
			_, err = syscall.Write(cli, []byte("123131"))
			if err != nil {
				t.Logf(err.Error())
				return
			}
			time.Sleep(time.Second)
			syscall.Close(cli)
		}()

		for {
			sock, addr, err := syscall.Accept(sk)
			if err != nil {
				return
			} else {
				t.Logf("addr:%v", addr)
				go handle(sock)

				time.Sleep(time.Second * 5)
				break
			}
		}
	})
}
