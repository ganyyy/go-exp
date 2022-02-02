package generic2

import (
	"testing"
)

func TestRangeConstraints(t *testing.T) {
	RangeConstraints(make([]MyNumber, 10))
	RangeConstraints(make([]MyString, 10))

	RangeConstraints1(make([]MyNumber, 10), make([]MyNumber, 10))
	RangeConstraints2(make([]MyNumber, 10), make([]MyString, 10))

	RangeConstraints3(make([]MyString, 10), make([]MyNumber, 10))
}
