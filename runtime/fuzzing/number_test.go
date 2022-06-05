package fuzzing

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberStringCompare(t *testing.T) {
	var testCases = [][]int{
		{1, 2},
		{3, 4},
	}

	for _, tc := range testCases {
		var p1, p2 = tc[0], tc[1]
		var s1, s2 = strconv.Itoa(p1), strconv.Itoa(p2)
		assert.Equal(t, p1 >= p2, NumberStringCompare(s1, s2))
	}

}

func FuzzNumberStringCompare(f *testing.F) {
	var testCases = [][]int{
		{1, 2},
		{3, 4},
	}

	for _, tc := range testCases {
		f.Add(tc[0], tc[1])
	}

	f.Fuzz(func(t *testing.T, p1, p2 int) {
		if p1 < 0 || p2 < 0 {
			return
		}
		var s1, s2 = strconv.Itoa(p1), strconv.Itoa(p2)
		assert.Equal(t, p1 > p2, NumberStringCompare(s1, s2))
	})
}
