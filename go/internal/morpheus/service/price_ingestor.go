package service

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

type IngestorService struct {
	Logger logging.Logger
}

func NewIngestorService(logger logging.Logger) *IngestorService {
	return &IngestorService{
		Logger: logger,
	}
}

func (s *IngestorService) IngestPrice(ctx context.Context, price *shared_domain.Price) error {
	s.Logger.Info("Morpheus Processing price",
		"symbol", price.Symbol,
		"price", price.Price,
		"currency", price.Currency)

	// Here you could add:
	// - Validations
	// - Database storage
	// - Dispatch to other systems
	// - Apply business rules

	return nil
}
