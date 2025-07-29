package locate

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Perms uint8

func (p Perms) String() string {
	var ret = [4]byte{'-', '-', '-', '-'}
	if p&PermsPrivate != 0 {
		ret[3] = 'p'
	}
	if p&PermsExecute != 0 {
		ret[2] = 'x'
	}
	if p&PermsWrite != 0 {
		ret[1] = 'w'
	}
	if p&PermsRead != 0 {
		ret[0] = 'r'
	}
	return string(ret[:])
}

func (p Perms) Is(op Perms) bool {
	return p&op == op
}

func ParsePerms(s string) Perms {
	var ret Perms
	for _, c := range s {
		switch c {
		case 'p':
			ret |= PermsPrivate
		case 'x':
			ret |= PermsExecute
		case 'w':
			ret |= PermsWrite
		case 'r':
			ret |= PermsRead
		default:
			continue
		}
	}
	return ret
}

const (
	PermsNone    Perms = 0
	PermsRead    Perms = 1 << (iota - 1) // 0001
	PermsWrite                           // 0010
	PermsExecute                         // 0100
	PermsPrivate                         // 1000
)

type PluginMapping struct {
	Start, End uint64
	Perms      Perms
	Offset     uint64
	Dev        string
	Inode      uint64
	Path       string
}

// String returns a string representation of the PluginMapping.
func (m *PluginMapping) String() string {
	var sb strings.Builder
	sb.WriteString("Start: 0x")
	sb.WriteString(strconv.FormatUint(m.Start, 16))
	sb.WriteString(", End: 0x")
	sb.WriteString(strconv.FormatUint(m.End, 16))
	sb.WriteString(", Perms: ")
	sb.WriteString(m.Perms.String())
	sb.WriteString(", Offset: 0x")
	sb.WriteString(strconv.FormatUint(m.Offset, 16))
	sb.WriteString(", Dev: ")
	sb.WriteString(m.Dev)
	sb.WriteString(", Inode: 0x")
	sb.WriteString(strconv.FormatUint(m.Inode, 10))
	sb.WriteString(", Path: ")
	sb.WriteString(m.Path)
	return sb.String()
}

// Size returns the size of the PluginMapping.
func (m *PluginMapping) Size() uint64 {
	return m.End - m.Start
}

func TryDecText(encFuncs []byte, selfPath string) error {

	file, err := os.Open("/proc/self/maps")
	if err != nil {
		return err
	}
	defer file.Close()

	var locateBytes = []byte(selfPath)

	var mappings []PluginMapping
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if !bytes.Contains(line, locateBytes) {
			continue
		}

		var start, end, offset uint64
		var perms, dev, path string
		var inode uint64
		if _, err := fmt.Sscanf(string(line), "%x-%x %s %x %s %d %s",
			&start, &end, &perms, &offset, &dev, &inode, &path); err != nil {
			continue
		}

		perm := ParsePerms(perms)
		if !perm.Is(PermsRead | PermsExecute) {
			continue // 只处理可读和可执行的映射, 也就是 .text 段
		}
		if end <= start {
			continue // 无效的映射
		}

		mappings = append(mappings, PluginMapping{
			Start:  start,
			End:    end,
			Perms:  ParsePerms(perms),
			Offset: offset,
			Dev:    dev,
			Inode:  inode,
			Path:   path,
		})
	}
	for _, mapping := range mappings {
		fmt.Println("Found mapping:", mapping.String())
	}

	for offset, size := range encFuncIter(encFuncs) {

		for _, mapping := range mappings {
			if offset < mapping.Offset {
				continue
			}
			if offset >= mapping.Offset+mapping.Size() {
				continue
			}
			off := offset - mapping.Offset + mapping.Start
			if off+size > mapping.End {
				continue
			}
			fmt.Printf("Decrypting function at offset %#x, size %#x\n", off, size)

			ciphertext := RawMemoryAccess(uintptr(off), int(size))
			copyToLocation(uintptr(off), decText(ciphertext))
		}
	}
	return nil
}

func copyToLocation(location uintptr, data []byte) {
	f := RawMemoryAccess(location, len(data))

	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_WRITE)
	copy(f, data[:])
	mprotectCrossPage(location, len(data), syscall.PROT_READ|syscall.PROT_EXEC)
}

func pageStart(ptr uintptr) uintptr {
	return ptr &^ (uintptr(syscall.Getpagesize() - 1))
}

func mprotectCrossPage(addr uintptr, length int, prot int) {
	pageSize := syscall.Getpagesize()
	for p := pageStart(addr); p < addr+uintptr(length); p += uintptr(pageSize) {
		page := RawMemoryAccess(p, pageSize)
		err := syscall.Mprotect(page, prot)
		if err != nil {
			panic(err)
		}
	}
}
