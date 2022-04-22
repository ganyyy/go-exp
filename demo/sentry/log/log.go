package log

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

func Errorf(format string, args ...interface{}) {
	sentry.CaptureMessage(fmt.Sprintf(format, args...))
	sentry.Flush(time.Second)
}

func Recover(format string, args ...interface{}) {
	if err := recover(); err != nil {
		var eventId = sentry.CaptureException(fmt.Errorf("panic:%v, reason:%s", err, fmt.Sprintf(format, args...)))
		sentry.Flush(time.Second)
		fmt.Printf("[%v] found error:%v\n", eventId, err)
	}
}
