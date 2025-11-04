package logging

import (
	"io"
	"log/slog"
	"os"
)

// FileLogger wraps slog.Logger for our application
type FileLogger struct {
	*slog.Logger
	file *os.File // Keep reference to close later
}

func (f *FileLogger) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *FileLogger) Info(msg string, keysAndValues ...any) {
	f.Logger.Info(msg, keysAndValues...)
}

func (f *FileLogger) Debug(msg string, keysAndValues ...any) {
	f.Logger.Debug(msg, keysAndValues...)
}

func (f *FileLogger) Error(msg string, keysAndValues ...any) {
	f.Logger.Error(msg, keysAndValues...)
}

func (f *FileLogger) Warn(msg string, keysAndValues ...any) {
	f.Logger.Warn(msg, keysAndValues...)
}

func NewFileLogger(filePath string, logLevel string) (*FileLogger, error) {
	var level slog.Level

	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var writer io.Writer = os.Stdout // Default to stdout
	var file *os.File

	// Only open file if path is provided
	if filePath != "" {
		var err error
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writer = file
	}

	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)

	return &FileLogger{
		Logger: logger,
		file:   file,
	}, nil
}
