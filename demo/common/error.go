package common

import (
	"fmt"
	"sync/atomic"
)

//go:generate stringer -type innerErrCode -output const_string.go
type innerErrCode uint32

func (i innerErrCode) Code() uint32 {
	return uint32(i)
}

func (i innerErrCode) Error() string {
	return i.String()
}

func (i innerErrCode) Equal(code IErrorCode) bool {
	if inner, ok := code.(innerErrCode); ok {
		return inner == i
	}
	return false
}

type errInner struct {
	code IErrorCode
	msg  interface{}
}

func (e errInner) Error() error {
	if e.code == ErrNo {
		return nil
	}
	return fmt.Errorf("code:%v, error:%v", e.code, e.msg)
}

//ErrorInfo 通用的对外错误信息结构
type ErrorInfo struct {
	inner atomic.Value // 多线程安全
}

func (e *ErrorInfo) Set(code IErrorCode, msg interface{}) {
	e.inner.Store(errInner{
		code: code,
		msg:  msg,
	})
}

func (e *ErrorInfo) Error() error {
	var code = e.inner.Load()
	if code == nil {
		return nil
	}
	if eCode, ok := code.(errInner); ok {
		return eCode.Error()
	}
	return nil
}
