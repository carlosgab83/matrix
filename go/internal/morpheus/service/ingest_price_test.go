package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/carlosgab83/matrix/go/internal/morpheus/domain"
	"github.com/carlosgab83/matrix/go/internal/morpheus/service"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/shared/mocks"
)

func TestIngestPriceTest(t *testing.T) {
	// Setup
	ctx, _ := context.WithCancel(context.Background())

	cfg := domain.Config{}

	price := shared_domain.Price{
		Symbol:    "ANY",
		Price:     1,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	// Mocks
	mockPriceRepository := mocks.NewPriceRepository(t)
	mockPriceRepository.On("InsertPrice", ctx, price).Return(nil)
	mockPublisher := mocks.NewPublisher(t)
	mockPublisher.On("NewDBPrice", ctx, price).Return(nil)

	// Real
	logger, _ := logging.NewLogger(cfg.CommonConfig)

	// Initilization
	ingestorService := service.NewIngestorService(ctx, logger, mockPriceRepository, mockPublisher)

	// Call
	ingestorService.IngestPrice(ctx, &price)

	// Assertions
	mockPriceRepository.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
