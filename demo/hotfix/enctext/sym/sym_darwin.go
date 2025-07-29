package sym

import (
	"bytes"
	"cmp"
	"debug/macho"
	"fmt"
	"slices"
	"strings"
)

func EncText(data []byte, fn EncFunc, prefix string) ([]byte, []byte, error) {
	// 不处理 Macho-Fat 文件

	f, err := macho.NewFile(bytes.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	if f.Symtab == nil {
		return nil, nil, fmt.Errorf("no symbol table found")
	}

	var sec *macho.Section
	var secIndex int
	for idx, s := range f.Sections {
		if s.Name == "__text" && s.Seg == "__TEXT" {
			sec = s
			secIndex = idx
			fmt.Println("Found .text section at index", idx, "with address", sec.Addr, "and size", sec.Size, "align", sec.Align, "offset", sec.Offset)
			break
		}
	}
	secIndex++ // secIndex is the index of the section in the file's sections slice
	if sec == nil {
		return nil, nil, fmt.Errorf("no (__TEXT.__text) section found")
	}
	if f.Symtab == nil {
		return nil, nil, fmt.Errorf("no symbol table found")
	}

	var allSyms = f.Symtab.Syms[:0]
	var encSyms []macho.Symbol
	for _, sym := range f.Symtab.Syms {
		if sym.Sect == 0 || int(sym.Sect) >= len(f.Sections) {
			continue
		}
		if int(sym.Sect) != secIndex {
			continue
		}
		allSyms = append(allSyms, sym)
		if !strings.HasPrefix(sym.Name, prefix) {
			continue
		}
		fmt.Println("Found function:", sym.Name)
		encSyms = append(encSyms, sym)
	}

	if len(encSyms) == 0 {
		fmt.Println("No functions found to encrypt with prefix:", prefix)
		return data, nil, nil
	}

	slices.SortFunc(allSyms, func(a, b macho.Symbol) int {
		return cmp.Compare(a.Value, b.Value)
	})
	slices.SortFunc(encSyms, func(a, b macho.Symbol) int {
		return cmp.Compare(a.Value, b.Value)
	})
	total := len(allSyms)

	var secEnd = sec.Addr + sec.Size
	var idx int
	var encFunc []byte
	for _, sym := range encSyms {
		for idx < total && allSyms[idx].Value < sym.Value {
			idx++
		}
		if idx >= total || allSyms[idx].Value != sym.Value {
			fmt.Printf("Symbol %s not found in all symbols, skipping\n", sym.Name)
			continue
		}
		var endValue = secEnd
		if idx+1 < total {
			endValue = allSyms[idx+1].Value
		}
		size := endValue - sym.Value
		if size > secEnd {
			fmt.Printf("Symbol %s has size %d, which exceeds section end %d, skipping\n", sym.Name, size, secEnd)
			continue
		}
		offset := sym.Value - sec.Addr + uint64(sec.Offset)
		if offset+size > uint64(len(data)) {
			fmt.Printf("Symbol %s has offset %d and size %d, which exceeds data length %d, skipping\n", sym.Name, offset, size, len(data))
			continue
		}
		fmt.Println("Symbol:", sym.Name, "Value:", sym.Value, "Size:", size, "Offset:", offset)
		if err := fn(offset, size); err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt function %s: %w", sym.Name, err)
		}
		encFunc = AppendEncFunc(encFunc, offset, size)
	}

	return data, encFunc, nil
}
