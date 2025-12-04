package notification

import (
	"context"

	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/tank/domain"
)

type Notifier interface {
	Notify(ctx context.Context, chatID string, payload string) error
	Register(context.Context)
	Close() error
}

func NewNotifier(cfg domain.Config, logger logging.Logger) (Notifier, error) {
	return NewTelegramNotifier(cfg, logger)
}
