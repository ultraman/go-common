package rest

// RateLimiter 限流接口
type RateLimiter interface {
}

func NewRateLimiter() RateLimiter {
	return nil
}

// WithRetry 重试接口
type WithRetry interface{}

func NewWithRetry(maxRetries int) WithRetry {
	return nil
}
