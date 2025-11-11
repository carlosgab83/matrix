package symbol_fetch

import (
	"context"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type SymbolFetcher interface {
	BTCUSDFetch(context.Context) (*shared_domain.Price, error)
	ETHUSDFetch(context.Context) (*shared_domain.Price, error)
}

type SymbolFetcherImpl struct {
	ETHUSDFetcher
	BTCUSDFetcher
}

func NewSymbolFetcher() SymbolFetcher {
	return &SymbolFetcherImpl{}
}
