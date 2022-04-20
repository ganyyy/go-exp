package log

import (
	"fmt"
	"log"
	"os"
)

const (
	_INFO  = "[INF]"
	_WARN  = "[WRN]"
	_ERROR = "[ERR]"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func toLog(prefix, format string, args ...interface{}) {
	logger.Output(3, fmt.Sprintf("%s %s", prefix, fmt.Sprintf(format, args...)))
}

func Infof(format string, args ...interface{}) {
	toLog(_INFO, format, args...)
}

func Warnf(format string, args ...interface{}) {
	toLog(_WARN, format, args...)
}

func Errorf(format string, args ...interface{}) {
	toLog(_ERROR, format, args...)
}
