package api_1_18

import (
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	var lock sync.Mutex

	lock.Lock()

	go func() {
		if !lock.TryLock() {
			log.Printf("lock error!")
		} else {
			log.Printf("lock success!")
		}
	}()

	time.Sleep(time.Second)

	lock.Unlock()

}

func TestStringsCut(t *testing.T) {
	var s = "1231231"
	var b, a, f = strings.Cut(s, "12")
	t.Logf("src:%v, cut:(%v,%v,%v)", s, b, a, f)
	var info, _ = debug.ReadBuildInfo()
	t.Logf("%+v", info)
}
