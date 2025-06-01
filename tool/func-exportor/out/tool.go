package test

import "plugin"

func mastBind[T any](p *plugin.Plugin, pt *T, name string) {
	if sym, err := p.Lookup(name); err == nil {
		if v, ok := sym.(*T); ok {
			*pt = *v
		} else {
			panic("type assertion failed for " + name)
		}
	} else {
		panic("lookup symbol " + name + " failed: " + err.Error())
	}
}
