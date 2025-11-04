package ingestion

import (
	"context"
	"time"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	matrix_proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
)

type GRPCPriceIngestorServer struct {
	matrix_proto.UnimplementedPriceIngestorServer
	IngestorService IngestorService
}

func NewGRPCPriceIngestorServer(ingestorService IngestorService) *GRPCPriceIngestorServer {
	return &GRPCPriceIngestorServer{
		IngestorService: ingestorService,
	}
}

func (s *GRPCPriceIngestorServer) IngestPrice(ctx context.Context, req *matrix_proto.PriceMessage) (*matrix_proto.IngestResponse, error) {
	// Convert from proto to domain (adapter's responsibility)
	price := &shared_domain.Price{
		Symbol:    req.Symbol,
		Price:     req.Price,
		Currency:  req.Currency,
		Timestamp: time.Unix(req.Timestamp, 0),
	}

	// Call the domain service
	err := s.IngestorService.IngestPrice(ctx, price)
	if err != nil {
		return &matrix_proto.IngestResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &matrix_proto.IngestResponse{
		Success: true,
		Message: "Price ingested successfully",
	}, nil
}
