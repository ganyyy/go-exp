package logger

import (
	"io"
	"log"
	"os"

	"google.golang.org/grpc/grpclog"
)

type logger struct {
	*log.Logger
}

func newLogger() io.Writer {
	return &logger{
		Logger: log.New(os.Stdout, "[GRPC]", 0),
	}
}

func (l *logger) Write(data []byte) (int, error) {
	l.Printf("%s", data)
	return len(data), nil
}

func SetGRPCLogger() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(newLogger(), newLogger(), newLogger()))
}
