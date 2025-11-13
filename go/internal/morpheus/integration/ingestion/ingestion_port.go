package ingestion

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type IngestorServiceInterface interface {
	IngestPrice(context.Context, *shared_domain.Price) error
}
