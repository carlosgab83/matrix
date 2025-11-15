package publication

import (
	"context"

	"github.com/carlosgab83/matrix/go/internal/morpheus/domain"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

const NewDBPriceTopic string = "price.db.new"

type Publicator interface {
	NewDBPrice(context.Context, shared_domain.Price) error
	Close() error
}

func NewPublicator(cfg domain.Config, logger logging.Logger) (Publicator, error) {
	return NewKafkaPublicator(cfg, logger)
}
