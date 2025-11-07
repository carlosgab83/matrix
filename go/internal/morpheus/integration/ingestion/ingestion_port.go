package ingestion

import (
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type IngestorServiceInterface interface {
	IngestPrice(price *shared_domain.Price) error
}
