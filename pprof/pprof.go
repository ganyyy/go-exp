package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
	"syscall"
)

type ProfileMonitor interface {
	Init(path string) ProfileMonitor
	Start() ProfileMonitor
	Done()
}

type ProfileBase struct {
	path string
	f    *os.File
}

func (p *ProfileBase) Init(path string) {
	p.path = path
	_ = syscall.Unlink(p.path)
	var f, err = os.Create(p.path)
	if err != nil {
		panic(fmt.Sprintf("Create %v error %v", p.path, err))
	}
	p.f = f
}

type CPUProfile struct {
	ProfileBase
}

func (c *CPUProfile) Init(path string) ProfileMonitor {
	c.ProfileBase.Init(path)
	return c
}

func (c *CPUProfile) Start() ProfileMonitor {
	var err = pprof.StartCPUProfile(c.f)
	if err != nil {
		panic("Start CPU Profile error :%v")
	}
	return c
}

func (c *CPUProfile) Done() {
	pprof.StopCPUProfile()
	_ = c.f.Close()
}

type MemProfile struct {
	ProfileBase
}

func (m *MemProfile) Init(path string) ProfileMonitor {
	m.ProfileBase.Init(path)
	return m
}

func (m *MemProfile) Start() ProfileMonitor {
	return m
}

func (m *MemProfile) Done() {
	_ = pprof.WriteHeapProfile(m.f)
}

type HTTPProfile struct {
	addr string
}

func (h *HTTPProfile) Init(addr string) ProfileMonitor {
	h.addr = addr
	return h
}

func (h *HTTPProfile) Start() ProfileMonitor {
	go func() {
		if err := http.ListenAndServe(h.addr, nil); err != nil {
			panic(err)
		}
	}()
	return h
}

func (h *HTTPProfile) Done() {
}
