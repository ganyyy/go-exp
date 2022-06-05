package main

import (
	"testing"
)

type MaybeNilInterface interface {
	Get() int
}

type EmptyGet struct{}

func (n EmptyGet) Get() int {
	return 0
}

type NilGet struct{}

func (n *NilGet) Get() int {
	return 0
}

func TestNilCheckInterface(t *testing.T) {
}
