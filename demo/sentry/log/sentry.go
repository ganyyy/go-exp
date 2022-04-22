package log

import (
	"github.com/getsentry/sentry-go"
)

/*

 */

func init() {
	sentry.Init(sentry.ClientOptions{
		Dsn:              "http://4d15e2eafc9749a8affa15afdc99a10f@localhost:9000/2",
		Debug:            false,
		AttachStacktrace: true,
		SampleRate:       0,
		TracesSampleRate: 0,
		TracesSampler:    nil,
		IgnoreErrors:     []string{},
		DebugWriter:      nil,
		Transport:        nil,
		ServerName:       "",
		Release:          "",
		Dist:             "",
	})
}
