package generic3

import "testing"

func TestLogger(t *testing.T) {
	DoLogger(&MyLogger{})
}
