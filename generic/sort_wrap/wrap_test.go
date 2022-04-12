package sortwrap

import (
	"testing"
)

func Test(t *testing.T) {
	var src = []int{
		1, 2, 3, 4, 5, 6, 9, 7,
	}

	Sort(src, func(t1, t2 int) bool {
		return t1 > t2
	})

	t.Logf("%+v", src)
}
