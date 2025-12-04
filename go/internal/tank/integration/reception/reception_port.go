package reception

import (
	"context"

	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/tank/domain"
)

type Receptor interface {
	BeginConsumption() error
	ReceiveCh() chan domain.NotificationPayload
	Close() error
}

func NewReceptor(ctx context.Context, cfg domain.Config, logger logging.Logger) (Receptor, error) {
	return NewKafkaReceptor(ctx, cfg, logger)
}
