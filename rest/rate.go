package rest

import "context"

// RateLimiter 限流接口
type RateLimiter interface {
	TryAccept() bool
	Accept()
	Stop()
	QPS() float32
	Wait(ctx context.Context) error
}

func NewRateLimiter() RateLimiter {
	return nil
}

// WithRetry 重试接口
type WithRetry interface {
	MaxRetries(max int)
	Retry()
}

func NewWithRetry(maxRetries int) WithRetry {
	return nil
}
