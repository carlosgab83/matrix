package service

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	shared_port "github.com/carlosgab83/matrix/go/internal/shared/port"
	"github.com/carlosgab83/matrix/go/internal/trinity/port"
)

type PriceIngestorService struct {
	logger shared_port.Logger
}

func NewPriceIngestorService(logger shared_port.Logger) port.PriceService {
	return &PriceIngestorService{
		logger: logger,
	}
}

func (s *PriceIngestorService) IngestPrice(ctx context.Context, price *shared_domain.Price) error {
	// Aquí va tu lógica de negocio
	s.logger.Info("Processing price",
		"symbol", price.Symbol,
		"price", price.Price,
		"currency", price.Currency)

	// Aquí podrías agregar:
	// - Validaciones
	// - Almacenamiento en BD
	// - Envío a otros sistemas
	// - Aplicar reglas de negocio

	return nil
}
