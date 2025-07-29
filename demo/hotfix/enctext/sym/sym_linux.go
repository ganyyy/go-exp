package sym

import (
	"bytes"
	"debug/elf"
	"fmt"
	"strings"
)

func EncText(data []byte, fn EncFunc, prefix string) ([]byte, []byte, error) {
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
		if err := fn(off, sym.Size); err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt function %s: %w", sym.Name, err)
		}
		encFunc = AppendEncFunc(encFunc, off, sym.Size)
	}
	return data, encFunc, nil
}
