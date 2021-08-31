package zerolog

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/yaoliu/go-common/logger"
	"os"
	"time"
)

type zeroLogger struct {
	zlogger zerolog.Logger
	options Options
}

var _ logger.Interface = &zeroLogger{}

func (l *zeroLogger) Init(opts ...logger.Option) error {
	zlog.Info().Msg("logger init start")
	for _, o := range opts {
		o(&l.options.Options)
	}

	if tf, ok := l.options.Context.Value(timeFieldFormatKey{}).(string); ok {
		l.options.TimeFieldFormat = tf
	}

	if mm, ok := l.options.Context.Value(messageFieldName{}).(string); ok {
		l.options.MessageFieldName = mm
	}

	// default
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "msg"
	zerolog.CallerSkipFrameCount = l.options.CallerSkipCount

	level := LevelToZeroLevel(l.options.Level)
	if level == 0 {
		level = zerolog.InfoLevel
	}

	l.zlogger = zerolog.New(l.options.Out).Level(level).With().Timestamp().Stack().Logger()
	if len(l.options.TimeFieldFormat) > 0 {
		zerolog.TimeFieldFormat = l.options.TimeFieldFormat
	}

	if len(l.options.MessageFieldName) > 0 {
		zerolog.MessageFieldName = l.options.MessageFieldName
	}

	zlog.Logger = l.zlogger
	return nil
}

func (l *zeroLogger) Info(ctx context.Context, message string, args ...interface{}) {
	l.zlogger.Info().Msgf(message, args...)
}

func (l *zeroLogger) Error(ctx context.Context, message string, args ...interface{}) {
	l.Logf(logger.ErrorLevel, message, args...)
}

func (l *zeroLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	l.Logf(logger.WarnLevel, message, args...)
}

func (l *zeroLogger) Debug(ctx context.Context, message string, args ...interface{}) {
	l.Logf(logger.DebugLevel, message, args...)
}

func (l *zeroLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	return
}

func (l *zeroLogger) Fatal(ctx context.Context, message string, args ...interface{}) {
	l.Logf(logger.FatalLevel, message, args...)
}

func (l *zeroLogger) Log(level logger.Level, args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.zlogger.WithLevel(LevelToZeroLevel(level)).Msg(msg)
	if level == logger.FatalLevel {
		os.Exit(1)
	}
}

func (l *zeroLogger) Logf(level logger.Level, message string, args ...interface{}) {
	l.zlogger.WithLevel(LevelToZeroLevel(level)).Msgf(message, args...)
	if level == logger.FatalLevel {
		os.Exit(1)
	}
}

func NewLogger(opts ...logger.Option) logger.Interface {
	options := Options{
		Options: logger.Options{
			Level:           logger.DebugLevel,
			Out:             logger.DefaultOutputer,
			CallerSkipCount: 4,
			Context:         context.Background(),
		},
		TimeFieldFormat:  time.RFC3339,
		MessageFieldName: "msg",
	}
	l := &zeroLogger{options: options}
	_ = l.Init(opts...)
	return l
}

func LevelToZeroLevel(level logger.Level) zerolog.Level {
	switch level {
	case logger.InfoLevel:
		return zerolog.InfoLevel
	case logger.DebugLevel:
		return zerolog.DebugLevel
	case logger.WarnLevel:
		return zerolog.WarnLevel
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	case logger.FatalLevel:
		return zerolog.FatalLevel
	case logger.TraceLevel:
		return zerolog.TraceLevel
	}
	return 0
}
