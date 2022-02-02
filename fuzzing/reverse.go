package fuzzing

import (
	"errors"
	"unicode/utf8"
)

var (
	ErrNotValidString = errors.New("not valid UTF-8 string")
)

func Reverse(s string) (string, error) {

	if !utf8.ValidString(s) {
		return s, ErrNotValidString
	}

	var bs = []rune(s)
	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}

	return string(bs), nil
}
