package zerolog

import "github.com/yaoliu/go-common/logger"

type Options struct {
	logger.Options
	MessageFieldName string
	TimeFieldFormat  string
}

type timeFieldFormatKey struct{}

func WithTimeFormat(timeFormat string) logger.Option {
	return logger.SetCustomOptions(timeFieldFormatKey{}, timeFormat)
}

type messageFieldName struct{}

func WithMessageFieldName(name string) logger.Option {
	return logger.SetCustomOptions(messageFieldName{}, name)
}
