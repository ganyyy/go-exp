// Code generated by "stringer -type=Data,Data2"; DO NOT EDIT.

package enum

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[D1-0]
	_ = x[D2-1]
	_ = x[D3-2]
	_ = x[D4-3]
	_ = x[D6-5]
	_ = x[D7-6]
	_ = x[D8-7]
}

const (
	_Data_name_0 = "D1D2D3D4"
	_Data_name_1 = "D6D7D8"
)

var (
	_Data_index_0 = [...]uint8{0, 2, 4, 6, 8}
	_Data_index_1 = [...]uint8{0, 2, 4, 6}
)

func (i Data) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _Data_name_0[_Data_index_0[i]:_Data_index_0[i+1]]
	case 5 <= i && i <= 7:
		i -= 5
		return _Data_name_1[_Data_index_1[i]:_Data_index_1[i+1]]
	default:
		return "Data(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

var _Data_stringToData = map[string]Data{
	D1.String(): D1,
	D2.String(): D2,
	D3.String(): D3,
	D4.String(): D4,
	D6.String(): D6,
	D7.String(): D7,
	D8.String(): D8,
}

func DataFromString(str string) (Data, bool) {
	val, ok := _Data_stringToData[str]
	return val, ok
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DD1-0]
	_ = x[DD2-1]
	_ = x[DD3-2]
	_ = x[DD4-3]
}

const _Data2_name = "DD1DD2DD3DD4"

var _Data2_index = [...]uint8{0, 3, 6, 9, 12}

func (i Data2) String() string {
	if i < 0 || i >= Data2(len(_Data2_index)-1) {
		return "Data2(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Data2_name[_Data2_index[i]:_Data2_index[i+1]]
}

var _Data2_stringToData = map[string]Data2{
	DD1.String(): DD1,
	DD2.String(): DD2,
	DD3.String(): DD3,
	DD4.String(): DD4,
}

func Data2FromString(str string) (Data2, bool) {
	val, ok := _Data2_stringToData[str]
	return val, ok
}