package zerolog

import (
	"context"
	"github.com/yaoliu/go-common/logger"
	"testing"
	"time"
)

func TestZeroLogger_Init(t *testing.T) {
	log := NewLogger(
		logger.WithLevel(logger.InfoLevel),
		logger.WithOutputer(logger.NewOutputer("test", "./")),
		logger.WithCallerSkipCount(2),
		WithTimeFormat(time.RFC3339))
	log.Info(context.TODO(), "%s hello world!", "yao")
}

func TestDefaultLogger(t *testing.T) {
	logger.DefaultLogger = NewLogger(
		logger.WithLevel(logger.InfoLevel),
		logger.WithOutputer(logger.NewOutputer("test", "./")),
		logger.WithCallerSkipCount(2),
		WithTimeFormat(time.RFC3339))
	logger.Info("%s hello world!", "yao")
}
