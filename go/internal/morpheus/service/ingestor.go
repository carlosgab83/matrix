package service

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

type IngestorService struct {
	Logger logging.Logger
	Ctx    context.Context
}

func NewIngestorService(ctx context.Context, logger logging.Logger) *IngestorService {
	return &IngestorService{
		Logger: logger,
		Ctx:    ctx,
	}
}

func (s *IngestorService) IngestPrice(price *shared_domain.Price) error {
	select {
	case <-s.Ctx.Done():
		s.Logger.Info("Morpheus stop due", "context", s.Ctx.Err())
		return s.Ctx.Err()
	default:
	}

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
