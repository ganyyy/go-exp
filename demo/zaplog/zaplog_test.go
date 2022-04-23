package log_test

import (
	"testing"
	"time"

	log "ganyyy.com/go-exp/demo/zaplog"
)

func TestZapLogger(t *testing.T) {
	var cfg log.Config
	cfg.FilePath = "./log/test_log.log"
	cfg.MaxAges = 10
	cfg.MaxSize = 1
	cfg.MaxBackups = 2
	cfg.Product = false
	log.Init(cfg)
	defer log.Sync()

	for i := 0; i < 10; i++ {
		log.Debugf("hello %s", "world")
		log.Infof("hello %s", "world")
		log.Errorf("hello %s", "world")
		time.Sleep(time.Second * 5)
	}
}
