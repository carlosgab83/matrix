package collector_test

import (
	"testing"
	"time"

	"github.com/carlosgab83/matrix/go/internal/neo/domain"
	"github.com/carlosgab83/matrix/go/internal/neo/service/collector"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/shared/mocks"
	"github.com/stretchr/testify/mock"
)

func TestNewCollector(t *testing.T) {
	// Setup
	price := shared_domain.Price{
		Symbol:    "ANY",
		Price:     1,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	cfg := domain.Config{
		DefaultFetchIntervalSeconds: 60,
		WorkersCount:                2,
		IngestorAddress:             "localhost:50051",
		Symbols: []domain.Symbol{
			{
				Nemo:                 "BTCUSD",
				Name:                 "Bitcoin",
				FetchIntervalSeconds: 1,
			},
			{
				Nemo:                 "ETHUSD",
				Name:                 "Ethreum",
				FetchIntervalSeconds: 1,
			},
		},
	}

	// Mocks
	mockIngestor := mocks.NewIngestor(t)
	mockIngestor.On("IngestPrice", mock.Anything, &price).Return(nil)
	mockSymbolFetcher := mocks.NewSymbolFetcher(t)
	mockSymbolFetcher.On("BTCUSDFetch", mock.Anything).Return(&price, nil)
	mockSymbolFetcher.On("ETHUSDFetch", mock.Anything).Return(&price, nil)

	// Real
	logger, _ := logging.NewLogger(cfg.CommonConfig)

	// Initilization
	col := collector.NewCollector(cfg, logger, mockIngestor, mockSymbolFetcher)

	// Call
	go col.Collect()
	time.Sleep(3 * time.Second)
	col.Stop()

	// Assertions
	mockSymbolFetcher.AssertExpectations(t)
	mockIngestor.AssertExpectations(t)
}
