package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestOutput(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	var csvWriter = csv.NewWriter(buf)
	defer func() {
		err := os.WriteFile("test.csv", buf.Bytes(), os.ModePerm)
		if err != nil {
			t.Logf("write error: %v", err)
		}
	}()
	defer csvWriter.Flush()

	for i := 0; i < 10; i++ {
		csvWriter.Write([]string{
			strconv.Itoa(i),
			strconv.Itoa(i) + ".data",
			strconv.Itoa(i) + ".file",
		})
	}
}
