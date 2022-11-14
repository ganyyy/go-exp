package test

import (
	"net/url"
	"redis-key-backup/config"
	"testing"

	"github.com/urfave/cli"
)

func TestConfig_Parse(t *testing.T) {
	var testCases = []string{
		"123:312",
		"1312:",
		"http://100.32.123.124:8899",
		"localhost:6379",
		"127.0.0.1:6379",
	}

	for _, str := range testCases {
		parseUrl, err := url.Parse(str)
		t.Logf("info:%+v, error:%v", parseUrl, err)
	}
}
func TestFlagLength(t *testing.T) {
	t.Logf("Len:%v, Cap:%v", len(config.RedisFlags), cap(config.RedisFlags))
	after := append(config.RedisFlags, cli.StringFlag{})
	t.Logf("Len:%v, Cap:%v", len(config.RedisFlags), cap(config.RedisFlags))
	t.Logf("Len:%v, Cap:%v", len(after), cap(after))
}
