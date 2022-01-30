package generic2

import (
	"testing"
)

func TestRangeConstraints(t *testing.T) {
	RangeConstraints(make([]MyNumber, 10))
	RangeConstraints(make([]MyString, 10))
}
