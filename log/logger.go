package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os/user"
	"time"
)

var LoggerRoot *zap.SugaredLogger
var parent *zap.SugaredLogger

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       true,
		DisableStacktrace: false,
		Encoding:          "console",
		OutputPaths:       []string{usr.HomeDir + "/iceapple.log"},
		ErrorOutputPaths:  []string{usr.HomeDir + "/iceapple_err.log"},
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
	LoggerRoot = parent.Named("[root]")
}

func NewLogger(name string) *zap.SugaredLogger {
	return parent.Named(name)
}

func Flush() {
	parent.Sync()
}
