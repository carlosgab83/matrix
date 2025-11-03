package platform

import (
	"io"
	"log/slog"
	"os"
)

// Logger wraps slog.Logger for our application
type Logger struct {
	*slog.Logger
	file *os.File // Keep reference to close later
}

// Close closes the underlying file (if any)
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

func NewLogger(logFilePath string, levelStr string) (*Logger, error) {
	var level slog.Level

	switch levelStr {
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
	if logFilePath != "" {
		var err error
		file, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writer = file
	}

	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
		file:   file,
	}, nil
}
