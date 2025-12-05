package publication

import (
	"context"

	"github.com/carlosgab83/matrix/go/internal/morpheus/domain"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

const NewDBPriceTopic string = "price.db.new"

type Publisher interface {
	NewDBPrice(context.Context, shared_domain.Price) error
	Close() error
}

func NewPublisher(cfg domain.Config, logger logging.Logger) (Publisher, error) {
	return NewKafkaPublisher(cfg, logger)
}
