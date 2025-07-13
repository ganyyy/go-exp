package main

import (
	"bytes"
	"crypto/rc4"
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"strings"
)

var Key = []byte("123456")

func main() {
	var input, prefix string
	flag.StringVar(&input, "input", "", "input dynamic library path")
	flag.StringVar(&prefix, "prefix", "", "prefix for function names")
	flag.Parse()

	if input == "" || prefix == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	f, err := elf.NewFile(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	var sec *elf.Section
	var idx int
	for i, s := range f.Sections {
		if s.Name == ".text" {
			sec = s
			idx = i
			break
		}
	}
	if sec == nil {
		panic("no .text section found")
	}
	fmt.Println("Found .text section at index", idx, "with address", sec.Addr, "and size", sec.Size, "align", sec.Addralign, "offset", sec.Offset)

	syms, err := f.Symbols()
	if err != nil {
		panic(err)
	}

	var cipher, _ = rc4.NewCipher(Key)

	var encFunc []byte

	for _, sym := range syms {
		if elf.ST_TYPE(sym.Info) != elf.STT_FUNC || int(sym.Section) != idx {
			continue
		}

		if !strings.HasPrefix(sym.Name, prefix) {
			continue
		}

		if sym.Value < sec.Addr || sym.Value+sym.Size > sec.Addr+sec.Size {
			continue
		}
		if sym.Size == 0 {
			continue
		}
		off := sym.Value - sec.Addr + sec.Offset
		if off+sym.Size > uint64(len(data)) {
			continue
		}
		fmt.Println("Symbol:", sym.Name, "Value:", sym.Value, "Size:", sym.Size, "Offset", off)
		enc := make([]byte, sym.Size)
		cipher.XORKeyStream(enc, data[off:off+sym.Size])
		copy(data[off:off+sym.Size], enc)

		encFunc = append(encFunc, fmt.Sprintf("%d:%d\n", off, sym.Size)...)
	}

	if err := os.WriteFile(input, data, 0644); err != nil {
		panic(err)
	}
	if err := os.WriteFile(input+".func", encFunc, 0644); err != nil {
		panic(err)
	}
	fmt.Println("Encryption complete, modified file saved.")
}
