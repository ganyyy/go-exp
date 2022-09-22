package main

import (
	"io"
	"strings"
	"testing"
)

func TestNumtiIO(t *testing.T) {
	// Create two readers
	foo := strings.NewReader("Hello Foo\n")
	bar := strings.NewReader("Hello Bar")

	// Create a multi reader
	mr := io.MultiReader(foo, bar)

	// Read data from multi reader
	for {

		var buf [3]byte
		n, err := mr.Read(buf[:])
		if err != nil {
			break
		}
		// Optional: Verify data
		t.Log(n, string(buf[:]))
	}

}
