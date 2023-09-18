package file

import (
	"context"
	"github.com/ultraman/go-common/logger"
	"time"
)

type fileLogger struct {
}

func (f fileLogger) Init(option ...logger.Option) error {
	panic("implement me")
}

func (f fileLogger) Log(level logger.Level, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Logf(level logger.Level, message string, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Info(ctx context.Context, message string, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Error(ctx context.Context, message string, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Debug(ctx context.Context, message string, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Fatal(ctx context.Context, message string, args ...interface{}) {
	panic("implement me")
}

func (f fileLogger) Trace(ctx context.Context, begin time.Time, fn func() (string, int64), err error) {
	panic("implement me")
}

var _ logger.Interface = &fileLogger{}
