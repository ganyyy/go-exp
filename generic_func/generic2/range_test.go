package generic2

import (
	"fmt"
	"strconv"
	"testing"
)

func TestRange(t *testing.T) {

	var src = make([]*Stu, 3)
	for i := range src {
		src[i] = &Stu{
			Name: strconv.Itoa(i),
			Age:  i,
		}
	}
	var src2 = make([]*Stu2, 3)
	for i := range src2 {
		if i&1 == 0 {
			continue
		}
		src2[i] = &Stu2{
			Stu: Stu{
				Name: strconv.Itoa(i + 10),
				Age:  i + 1,
			},
			Other: strconv.Itoa(i + 20),
		}
	}
	Range(src) // 泛型的优点
	Range2(src)
	Range(src2)

	var tmpSrc = make([]fmt.Stringer, 0, len(src)+len(src2))
	for i := range src {
		tmpSrc = append(tmpSrc, src[i])
	}
	for i := range src2 {
		tmpSrc = append(tmpSrc, src2[i])
	}

	RangeInterface(tmpSrc)
	Range(tmpSrc) // 完全兼容接口迭代
}
