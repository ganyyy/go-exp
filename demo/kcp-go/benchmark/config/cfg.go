package kcp_benchmark_config

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"

	"net/http"
	_ "net/http/pprof"

	_ "ganyyy.com/go-exp/helper"
	"github.com/BurntSushi/toml"
)

var Config struct {
	ServerAddr string // 服务端地址
	ListenNum  int    // 如果是服务端的话, 并行监听的数量

	ClientNum    int // 如果是客户端的话, 客户端数量
	EchoInterval int // 如果是客户端的话, echo间隔(毫秒)
	Intervals    int // 间隔时间
	PProfPort    int // pprof 端口

	IsKCP      bool // 是否是kcp
	IsServer   bool // 是否是服务端
	WriteDelay bool // 是否延迟写
}

var Order = binary.BigEndian

func ReadConfig(path string) error {
	_, err := toml.DecodeFile(path, &Config)
	return err
}

func LogAndExit(err error) {
	slog.Error("exit", slog.String("err", err.Error()))
	os.Exit(1)
}

func MustReadConfig(path string) {
	if err := ReadConfig(path); err != nil {
		LogAndExit(fmt.Errorf("read config error: %v", err))
	}
}

func OpenPProf() {
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", Config.PProfPort), nil)
		if err != nil {
			LogAndExit(fmt.Errorf("start pprof error: %v", err))
		}
	}()
}
