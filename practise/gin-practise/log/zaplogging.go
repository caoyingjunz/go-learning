package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func NewZapLogger(c Configuration) (Logger, error) {
	var w io.Writer
	switch c.LogFile {
	case "stdout":
		w = os.Stdout
	case "stderr":
		w = os.Stderr
	default:
		w = &lumberjack.Logger{
			Filename:   c.LogFile,
			MaxSize:    c.RotateMaxSize,
			MaxAge:     c.RotateMaxAge,
			MaxBackups: c.RotateMaxBackups,
			Compress:   c.Compress,
		}
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(c.LogLevel)); err != nil {
		return nil, err
	}

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.NewMultiWriteSyncer(zapcore.AddSync(w)), level)

	var cores []zapcore.Core
	cores = append(cores, core)
	Tee := zapcore.NewTee(cores...)
	logger := zap.New(Tee, zap.AddCaller(), zap.AddCallerSkip(1))
	return &zapLogger{
		logger:    logger,
		writer:    w,
		verbosity: 0,
	}, nil
}

type zapLogger struct {
	logger    *zap.Logger
	writer    io.Writer
	verbosity int
}

func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *zapLogger) Infof(f string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(f, args...))
}
