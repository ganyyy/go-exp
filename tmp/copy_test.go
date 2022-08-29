package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"
)

type ICopy interface {
	F1()
	F2()
}

type CopyBase struct {
	Wg sync.WaitGroup
}

func (c *CopyBase) F1() {
	log.Printf("Base F1")
}

func (c *CopyBase) F2() {
	log.Printf("Base F2")
}

type Copy1 struct {
	CopyBase
}

func (c *Copy1) F1() {
	c.Wg.Add(1)
	log.Printf("Copy1 F1")
}

func (c *Copy1) F2() {
	c.Wg.Done()
	log.Printf("Copy1 F2")
}

func TestCopy(t *testing.T) {
	var c1 Copy1
	var c ICopy = &c1
	c.F1()
	c.F2()
	c1.Wg.Wait()

	_ = fmt.Errorf("%s %w", "123", errors.New("123"))
}
