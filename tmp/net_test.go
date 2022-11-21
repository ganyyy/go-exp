package main

import (
	"net"
	"testing"
	"time"
)

func TestNetConnect(t *testing.T) {
	// 系统套接字创建时会设置为非阻塞的, 这时候就可以通过context控制连接超时
	_, err := net.DialTimeout("tcp", "localhost:6666", time.Second)
	t.Logf("conn error:%v", err)
}
