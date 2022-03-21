package inter

import (
	"testing"
)

func TestLocker(t *testing.T) {

}

type Runner interface {
	Do()
}

type Run struct{ _ [10]string }

func (r *Run) Do() {}

type Run2 struct{ _ [10]string }

func (r Run2) Do() {}

func Benchmark(b *testing.B) {

	b.Run("direct pointer", func(b *testing.B) {
		var run = &Run{}
		for i := 0; i < b.N; i++ {
			run.Do()
		}
	})

	b.Run("direct struct", func(b *testing.B) {
		var run = Run2{}
		for i := 0; i < b.N; i++ {
			run.Do()
		}
	})

	b.Run("pointer", func(b *testing.B) {
		var run Runner = &Run{}
		for i := 0; i < b.N; i++ {
			run.Do()
		}
	})

	b.Run("empty", func(b *testing.B) {
		var run Runner = Run2{}
		for i := 0; i < b.N; i++ {
			run.Do()
		}
	})

	b.Run("pointer eface", func(b *testing.B) {
		var run any = &Run{}
		for i := 0; i < b.N; i++ {
			run.(Runner).Do()
		}
	})

	b.Run("struct eface", func(b *testing.B) {
		var run any = Run2{}
		for i := 0; i < b.N; i++ {
			run.(Runner).Do()
		}
	})

	b.Run("pointer conversion", func(b *testing.B) {
		var run any = &Run{}
		for i := 0; i < b.N; i++ {
			switch v := run.(type) {
			case Runner:
				v.Do()
			}
		}
	})

	b.Run("struct conversion", func(b *testing.B) {
		var run any = Run2{}
		for i := 0; i < b.N; i++ {
			switch v := run.(type) {
			case Runner:
				v.Do()
			}
		}
	})
}
