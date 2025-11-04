package logging

import (
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Close() error
}

func NewLogger(cfg shared_domain.CommonConfig) (Logger, error) {
	return NewFileLogger(cfg.LogFilePath, cfg.LogLevel)
}
