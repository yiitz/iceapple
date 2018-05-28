package log

import (
	"github.com/yiitz/iceapple/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"flag"
)

var LoggerRoot *zap.SugaredLogger
var parent *zap.SugaredLogger

var logLevel zap.AtomicLevel

func init() {
	outputPath := []string{storage.AppDir() + "/std.log"}
	if flag.Lookup("test.v") != nil {
		outputPath = append(outputPath, "stdout")
	}
	logLevel.UnmarshalText([]byte("debug"))
	cfg := zap.Config{
		Level:             logLevel,
		Development:       true,
		DisableStacktrace: false,
		Encoding:          "console",
		OutputPaths:       outputPath,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			TimeKey:       "time",
			NameKey:       "name",
			CallerKey:     "caller",
			LevelKey:      "level",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeName:    zapcore.FullNameEncoder,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			StacktraceKey: "stacktrace",
		},
	}

	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	parent = l.Sugar()
	LoggerRoot = NewLogger("root")
}

func SetLevel(level string) {
	logLevel.UnmarshalText([]byte(level))
}

func NewLogger(name string) *zap.SugaredLogger {
	return parent.Named("[" + name + "]")
}

func Flush() {
	parent.Sync()
}
