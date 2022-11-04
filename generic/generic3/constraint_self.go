package generic3

import (
	"log"
	"strings"
)

type Logger[T any] interface {
	WithField(string, string) T
	Info(string)
}

type MyLogger struct {
	fields []string
}

func (m *MyLogger) WithField(name, value string) *MyLogger {
	m.fields = append(m.fields, name+"="+value)
	return m
}

func (m *MyLogger) Info(s string) {
	log.Printf("[%s] %s", s, strings.Join(m.fields, ","))
}

// DoLogger 这是一个自引用接口, 可以做到完全的鸭子类型. 不需要任何侵入式导入
func DoLogger[T Logger[T]](t T) {
	t.WithField("1.18", "generic").Info("very good!")
}
