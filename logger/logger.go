package logger

import (
	"context"
	"fmt"
	"log"
	"time"
)

var (
	DefaultLogger = NewLogger()
)

func NewLogger(opts ...Option) Interface {
	return newDefaultLogger()
}

type Level int8

const (
	TraceLevel Level = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	}
	return ""
}

func GetLevel(level string) Level {
	switch level {
	case InfoLevel.String():
		return InfoLevel
	case ErrorLevel.String():
		return ErrorLevel
	case DebugLevel.String():
		return DebugLevel
	case TraceLevel.String():
		return TraceLevel
	case FatalLevel.String():
		return FatalLevel
	}
	return 0
}

type Interface interface {
	Init(...Option) error
	Log(level Level, args ...interface{})
	Logf(level Level, message string, args ...interface{})
	Info(ctx context.Context, message string, args ...interface{})
	Error(ctx context.Context, message string, args ...interface{})
	Warn(ctx context.Context, message string, args ...interface{})
	Debug(ctx context.Context, message string, args ...interface{})
	Fatal(ctx context.Context, message string, args ...interface{})
	Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error)
}

func Info(args ...interface{}) {
	DefaultLogger.Log(InfoLevel, args...)
}

func Infof(message string, args ...interface{}) {
	DefaultLogger.Logf(InfoLevel, message, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Log(WarnLevel, args...)
}

func Warnf(message string, args ...interface{}) {
	DefaultLogger.Logf(WarnLevel, message, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Log(ErrorLevel, args...)
}

func Errorf(message string, args ...interface{}) {
	DefaultLogger.Logf(ErrorLevel, message, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Log(DebugLevel, args...)
}

func Debugf(message string, args ...interface{}) {
	DefaultLogger.Logf(DebugLevel, message, args...)
}
func Fatal(args ...interface{}) {
	DefaultLogger.Log(FatalLevel, args...)
}

func Fatalf(message string, args ...interface{}) {
	DefaultLogger.Logf(FatalLevel, message, args...)
}

type defaultLogger struct{}

var _ Interface = &defaultLogger{}

func (d defaultLogger) Init(option ...Option) error {
	return nil
}

func (d defaultLogger) Log(level Level, args ...interface{}) {
	log.Println(args...)
}

func (d defaultLogger) Logf(level Level, message string, args ...interface{}) {
	log.Printf(message, args...)
}

func (d defaultLogger) Info(ctx context.Context, message string, args ...interface{}) {
	log.Printf(message, args...)
}

func (d defaultLogger) Error(ctx context.Context, message string, args ...interface{}) {
	log.Printf(message, args...)
}

func (d defaultLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	log.Printf(message, args...)
}

func (d defaultLogger) Debug(ctx context.Context, message string, args ...interface{}) {
	log.Printf(message, args...)
}

func (d defaultLogger) Fatal(ctx context.Context, message string, args ...interface{}) {
	log.Fatal(fmt.Sprintf(message, args...))
}

func (d defaultLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	return
}

func newDefaultLogger() Interface {
	return &defaultLogger{}
}
