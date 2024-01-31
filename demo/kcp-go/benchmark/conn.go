package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
	"github.com/xtaci/kcp-go/v5"
)

// EchoDataSize 8字节的纳秒级时间戳
const EchoDataSize = 8

var (
	ErrMsgSize = errors.New("invalid msg size")
)

type Conn struct {
	*slog.Logger
	net.Conn
	close     chan struct{}
	ticker    chan struct{}
	isClosed  atomic.Bool
	closeOnce sync.Once
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Logger: slog.Default().With(
			slog.String("type", "conn"),
			slog.String("remote", conn.RemoteAddr().String()),
			slog.String("local", conn.LocalAddr().String()),
		),
		Conn:   conn,
		close:  make(chan struct{}),
		ticker: make(chan struct{}),
	}
}

func (conn *Conn) Tick() bool {
	if conn.isClosed.Load() {
		return false
	}
	select {
	case conn.ticker <- struct{}{}:
	default:
	}
	return true
}

func (conn *Conn) Write(t int64) error {
	var buf [EchoDataSize]byte
	kcp_benchmark_config.Order.PutUint64(buf[:], uint64(t))
	n, err := conn.Conn.Write(buf[:])
	if n != EchoDataSize {
		return fmt.Errorf("invalid write size: %d", n)
	}
	return err
}

func (conn *Conn) Read() (int64, error) {
	var buf [EchoDataSize]byte
	n, err := io.ReadFull(conn.Conn, buf[:])
	if n != EchoDataSize {
		return 0, ErrMsgSize
	}
	if err != nil {
		return 0, err
	}
	return int64(kcp_benchmark_config.Order.Uint64(buf[:])), err
}

func (conn *Conn) StartWrite() {
	for {
		select {
		case <-conn.ticker:
			err := conn.Write(time.Now().UnixNano())
			if err != nil {
				conn.Logger.Error("write", slog.String("err", err.Error()))
				return
			}
		case <-conn.close:
			conn.Info("close")
			return
		}
	}
}

func (conn *Conn) StartRead() {
	defer conn.Close()
	for {
		t, err := conn.Read()
		// 记录往返延迟埋点
		sub := time.Duration(time.Now().UnixNano() - t)
		// conn.Logger.Info("echo", slog.String("sub", sub.String()))
		echoTimer.Update(sub)
		if err != nil {
			conn.Logger.Error("read", slog.String("err", err.Error()))
			return
		}
	}
}

func (conn *Conn) Close() {
	conn.closeOnce.Do(func() {
		conn.isClosed.Store(true)
		err := conn.Conn.Close()
		if err != nil {
			conn.Logger.Error("close", slog.String("err", err.Error()))
		}
		if conn.close != nil {
			close(conn.close)
			conn.close = nil
		}
	})
}

type AcceptCallback func(conn net.Conn) error

type Listener struct {
	net.Listener
	*slog.Logger
	AcceptCallback
}

func NewListener(listener net.Listener, afterAccept AcceptCallback) *Listener {
	return &Listener{
		Listener: listener,
		Logger: slog.Default().With(
			slog.String("type", "listener"),
			slog.String("local", listener.Addr().String()),
		),
		AcceptCallback: afterAccept,
	}
}

func (lis *Listener) Start() {
	lis.Info("start")
	for {
		conn, err := lis.Accept()
		if err != nil {
			lis.Error("accept", slog.String("err", err.Error()))
			return
		}
		// lis.Info("accept", slog.String("remote", conn.RemoteAddr().String()))
		go lis.HandleConn(conn)
	}
}

func (lis *Listener) HandleConn(conn net.Conn) {
	if lis.AcceptCallback != nil {
		err := lis.AcceptCallback(conn)
		if err != nil {
			lis.Error("after accept", slog.String("err", err.Error()))
			return
		}
	}
	c := NewConn(conn)
	defer c.Close()
	for {
		t, err := c.Read()
		if err != nil {
			c.Logger.Error("read", slog.String("err", err.Error()))
			return
		}
		err = c.Write(t)
		if err != nil {
			c.Logger.Error("write", slog.String("err", err.Error()))
			return
		}
	}
}

func AcceptKcpSession(conn net.Conn) error {
	c, ok := conn.(*kcp.UDPSession)
	if !ok {
		return nil
	}
	kcp_benchmark_config.InitKcpSession(c)
	return nil
}

func RunClients(num int, dial func(string) (net.Conn, error)) {
	slog.Info("RunClients", slog.Int("num", num))
	var allConns = make([]func() bool, 0, num)
	// dial
	for idx := 0; idx < num; idx++ {
		netConn, err := dial(kcp_benchmark_config.Config.ServerAddr)
		if err != nil {
			slog.Error("dial error", slog.Int("idx", idx), slog.String("err", err.Error()))
			continue
		}
		if kcpConn, ok := netConn.(*kcp.UDPSession); ok {
			kcp_benchmark_config.InitKcpSession(kcpConn)
		}
		conn := NewConn(netConn)
		go conn.StartWrite()
		go conn.StartRead()
		allConns = append(allConns, conn.Tick)
	}
	// ticker
	var ticker = time.NewTicker(time.Duration(kcp_benchmark_config.Config.EchoInterval) * time.Millisecond)
	defer ticker.Stop()

	time.Sleep(time.Second * 5)
	slog.Info("start ticker")
	for range ticker.C {
		ln := len(allConns)
		for idx := 0; idx < ln; {
			if allConns[idx] != nil && allConns[idx]() {
				idx++
			} else {
				allConns[idx] = allConns[ln-1]
				allConns[ln-1] = nil
				allConns = allConns[:ln-1]
				ln--
			}
		}
	}
}

func DailKCP(addr string) (net.Conn, error) {
	return kcp.Dial(addr)
}

func DailTCP(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}
