package persisence

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type PriceRepository interface {
	InsertPrice(context.Context, shared_domain.Price) error
	Close() error
}

func NewPriceRepository(connStr string) (PriceRepository, error) {
	return NewPgPriceRepository(connStr)
}
