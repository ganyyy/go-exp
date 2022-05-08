package patch

import "C"

import (
	"sync"
	"unsafe"
)

type funcVal struct {
	_   uintptr
	ptr unsafe.Pointer
}

func toCString(str string) *C.char {
	if len(str) == 0 {
		return (*C.char)(NULL)
	}
	var name = append([]byte(str), 0) // 增加一个\0终止符
	return (*C.char)(unsafe.Pointer(&name[0]))
}

var NULL = unsafe.Pointer((*int)(nil))

//TODO fix 多线程不安全

type Symbol struct {
	pointer unsafe.Pointer
	back    []byte
}

func (s *Symbol) Patch(symbol *Symbol) {
	if symbol == nil {
		return
	}
	if len(s.back) == 0 {
		s.back = Backup(s.pointer)
	}
	Patch(s.pointer, symbol.pointer)
}

func (s *Symbol) Restore() {
	if s.pointer == nil || len(s.back) == 0 {
		return
	}
	s.back = nil
	Restore(s.pointer, s.back)
}

type Plugin struct {
	handler unsafe.Pointer
	symbol  sync.Map
}

func NewPlugin(path string) (*Plugin, error) {
	var handler, err = PluginOpen(path)
	if err != nil {
		return nil, err
	}
	return &Plugin{
		handler: handler,
	}, err
}

func (p *Plugin) Symbol(name string) (*Symbol, error) {
	if symbol, ok := p.symbol.Load(name); ok {
		return symbol.(*Symbol), nil
	}
	var symbol, err = LookupSymbol(p.handler, name)
	if err != nil {
		return nil, err
	}
	var s = &Symbol{
		pointer: symbol,
	}
	p.symbol.Store(name, s)
	return s, nil
}

func (p *Plugin) Close() error {
	p.symbol.Range(func(key, value any) bool {
		var symbol, ok = value.(*Symbol)
		if !ok {
			return true
		}
		symbol.Restore()
		p.symbol.Delete(key)
		return true
	})

	return PluginClose(p.handler)
}
