package kcp_benchmark_config

import (
	"reflect"
	"unsafe"

	"github.com/xtaci/kcp-go/v5"
)

var (
	kcpOffsetOfUDPSession uintptr = 0
	rxMinRtoOffsetOfKcp   uintptr = 0
)

func init() {
	const (
		kcpFieldName      = "kcp"
		rxMinRtoFieldName = "rx_minrto"
	)

	tp := reflect.TypeOf(&kcp.UDPSession{}).Elem()
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		if field.Name == kcpFieldName {
			kcpOffsetOfUDPSession = field.Offset
			break
		}
	}
	tp = reflect.TypeOf(&kcp.KCP{}).Elem()
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		if field.Name == rxMinRtoFieldName {
			rxMinRtoOffsetOfKcp = field.Offset
			break
		}
	}
}

func SetMinRxRto(session *kcp.UDPSession, minRto uint32) {
	kcpPtr := (*kcp.KCP)(unsafe.Pointer(uintptr(unsafe.Pointer(session)) + kcpOffsetOfUDPSession))
	*(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(kcpPtr)) + rxMinRtoOffsetOfKcp)) = minRto
}

func InitKcpSession(session *kcp.UDPSession) {
	session.SetWriteDelay(Config.WriteDelay)
	session.SetNoDelay(1, Config.Intervals, 2, 1)
}
