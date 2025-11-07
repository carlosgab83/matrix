package ingestion

import (
	"fmt"
	"io"
	"time"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	matrix_proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
)

type GRPCPriceIngestorServer struct {
	matrix_proto.UnimplementedPriceIngestorServer
	IngestorService IngestorServiceInterface
	Logger          logging.Logger
}

func NewGRPCPriceIngestorServer(ingestorService IngestorServiceInterface, logger logging.Logger) *GRPCPriceIngestorServer {
	return &GRPCPriceIngestorServer{
		IngestorService: ingestorService,
		Logger:          logger,
	}
}

func (s *GRPCPriceIngestorServer) IngestPrice(stream matrix_proto.PriceIngestor_IngestPriceServer) error {
	// Convert from proto to domain (adapter's responsibility)
	for {
		priceMsg, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&matrix_proto.IngestResponse{
				Success: true,
			})

			return nil
		}

		if err != nil {
			return fmt.Errorf("receiving price error: %w", err)
		}

		price := &shared_domain.Price{
			Symbol:    priceMsg.Symbol,
			Price:     priceMsg.Price,
			Currency:  priceMsg.Currency,
			Timestamp: time.Unix(priceMsg.Timestamp, 0),
		}

		// Call the domain service
		err = s.IngestorService.IngestPrice(price)
		if err != nil {
			err = fmt.Errorf("ingesting price %v error: %w", price, err)
			s.Logger.Error("error ingesting price", "error", err)
		}
	}
}
