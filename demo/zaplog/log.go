package log

func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	//TODO 接入sentry
	sugar.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	sugar.Panicf(format, args...)
}
