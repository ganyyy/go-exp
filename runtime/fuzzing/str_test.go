package fuzzing

import "testing"

func FuzzStrCompare(f *testing.F) {
	var testCases = [][]string{
		{"123", "24"},
		{"12", "5"},
	}

	for _, tc := range testCases {
		f.Add(tc[0], tc[1])
	}

	f.Fuzz(func(t *testing.T, p1, p2 string) {

		var v1 = CompareString(p1, p2)
		var v2 = CompareString(p2, p1)
		if p1 == p2 {
			if v1 != v2 {
				t.Errorf("p1:%v, p2:%v compare not equal!", p1, p2)
			}
		} else if v1 != !v2 {
			t.Errorf("p1:%v, p2:%v, compare not valid!", p1, p2)
		}
	})
}
