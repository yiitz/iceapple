package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"github.com/yiitz/iceapple/storage"
)

var LoggerRoot *zap.SugaredLogger
var parent *zap.SugaredLogger

func init() {
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableStacktrace: false,
		Encoding:          "console",
		OutputPaths:       []string{storage.AppDir() + "/std.log"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			TimeKey:       "time",
			NameKey:       "name",
			CallerKey:     "caller",
			LevelKey:      "level",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeName:    zapcore.FullNameEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
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

func NewLogger(name string) *zap.SugaredLogger {
	return parent.Named("[" + name + "]")
}

func Flush() {
	parent.Sync()
}
