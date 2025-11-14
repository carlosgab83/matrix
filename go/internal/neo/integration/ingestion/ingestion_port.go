package ingestion

import (
	"context"

	"github.com/carlosgab83/matrix/go/internal/neo/domain"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

// Ingestor defines the interface for ingesting prices
type Ingestor interface {
	IngestPrice(context.Context, *shared_domain.Price) error
	Close() error
}

func NewIngestor(cfg domain.Config) (Ingestor, error) {
	return NewGRPCPriceIngestor(cfg.IngestorAddress, cfg.GRPCSharedToken)
}
