package logger

import "context"

type Options struct {
	Level           Level
	Out             Outputer
	CallerSkipCount int
	Context         context.Context
}

type Option func(*Options)

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}
func WithOutputer(out Outputer) Option {
	return func(o *Options) {
		o.Out = out
	}
}

func WithCallerSkipCount(c int) Option {
	return func(o *Options) {
		o.CallerSkipCount = c
	}
}

func SetCustomOptions(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
