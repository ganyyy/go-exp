package conn

import "sync"

var peer sync.Map

var val struct{}

func Add(addr Addr) {
	peer.LoadOrStore(addr, val)
}

func Del(addr Addr) {
	peer.Delete(addr)
}

func Has(addr Addr) bool {
	_, loaded := peer.Load(addr)
	return loaded
}

func List() []Addr {
	var ret []Addr
	peer.Range(func(addr, _ any) bool {
		ret = append(ret, addr.(Addr))
		return true
	})
	return ret
}
