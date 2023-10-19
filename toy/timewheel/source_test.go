package timewheel_test

import (
	"math"
	"runtime"
	"testing"
)

func panicIfError(n *Data) {
	if n.A < 0 {
		panic("error")
	}
}

type Data struct {
	A int
}

func Test123(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			var _stack [10]uintptr
			n := runtime.Callers(0, _stack[:])
			stack := _stack[:n]
			for i := 0; i < len(stack); i++ {
				f := runtime.FuncForPC(stack[i])
				file, line := f.FileLine(stack[i])
				t.Logf("%x %s:%d %s\n", stack[i], file, line, f.Name())
			}
		}
	}()
	panicIfError(nil)

}

func TestOOM(t *testing.T) {
	var mem = make([]byte, math.MaxInt32<<8)

	_ = mem

}
