package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

type Config struct {
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAges    int
	Product    bool
}

func init() {
	logger, _ = zap.NewDevelopment()
}

func getLoggerWriter(cfg Config) zapcore.WriteSyncer {
	var writer = &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAges,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  true,
	}

	return zapcore.AddSync(writer)
}

func Init(cfg Config) {
	var logLevel zapcore.Level
	if cfg.Product {
		logLevel = zapcore.InfoLevel
	} else {
		logLevel = zapcore.DebugLevel
	}

	var encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewTee(
		// 写入文件
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), getLoggerWriter(cfg), logLevel),
		// 输出到控制台
		// zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.Lock(os.Stdout), logLevel),
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // 再包装一层, 所以需要加一下
	sugar = logger.Sugar()
}

func Sync() {
	if logger == nil {
		return
	}
	_ = logger.Sync()
}
