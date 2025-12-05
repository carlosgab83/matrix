package service

import (
	"context"

	"github.com/carlosgab83/matrix/go/internal/morpheus/integration/persisence"
	"github.com/carlosgab83/matrix/go/internal/morpheus/integration/publication"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

type IngestorService struct {
	Logger          logging.Logger
	PriceRepository persisence.PriceRepository
	Publisher       publication.Publisher
	Ctx             context.Context
}

func NewIngestorService(ctx context.Context, logger logging.Logger, priceRepository persisence.PriceRepository, publisher publication.Publisher) *IngestorService {
	return &IngestorService{
		Logger:          logger,
		PriceRepository: priceRepository,
		Publisher:       publisher,
		Ctx:             ctx,
	}
}

func (s *IngestorService) IngestPrice(ctx context.Context, price *shared_domain.Price) error {
	select {
	case <-s.Ctx.Done():
		s.Logger.Info("Morpheus stop due", "context", s.Ctx.Err())
		return s.Ctx.Err()
	default:
	}

	err := s.PriceRepository.InsertPrice(s.Ctx, *price)
	if err != nil {
		s.Logger.Error("Error Inserting Price",
			"symbol", price.Symbol,
			"price", price.Price,
			"currency", price.Currency,
			"error", err)

		return nil
	}

	s.Logger.Info("Price inserted",
		"symbol", price.Symbol,
		"price", price.Price,
		"currency", price.Currency)

	s.Publisher.NewDBPrice(s.Ctx, *price)

	return nil
}
