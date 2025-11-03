package port

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type PriceService interface {
	IngestPrice(ctx context.Context, price *shared_domain.Price) error
}
