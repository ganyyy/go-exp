package log

import (
	"github.com/getsentry/sentry-go"
)

/*
!5xkwvz=3-d4f+_^@ow4)eguq5=duhbm=#ppu=h5eusesdtt)p
*/

func InitSentry() {
	sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Debug:            false,
		AttachStacktrace: false,
		SampleRate:       0,
		TracesSampleRate: 0,
		TracesSampler:    nil,
		DebugWriter:      nil,
		Transport:        nil,
		ServerName:       "",
		Release:          "",
		Dist:             "",
		Environment:      "",
		MaxBreadcrumbs:   0,
	})

}
