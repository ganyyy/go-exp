package log_test

import (
	"testing"

	"p2p/log"
)

func TestLogger(t *testing.T) {
	log.Errorf("%v, %v", 100, 200)
	log.Warnf("%v, %v", 200, 300)
	log.Infof("%v, %v", 200, 300)
}
