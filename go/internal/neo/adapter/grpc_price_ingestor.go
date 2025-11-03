package adapter

import (
	"context"
	"fmt"

	"github.com/carlosgab83/matrix/go/internal/neo/port"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCPriceIngestor implements the PriceIngestor port using gRPC
type GRPCPriceIngestor struct {
	client proto.PriceIngestorClient
	conn   *grpc.ClientConn
}

// NewGRPCPriceIngestor creates a new gRPC adapter
func NewGRPCPriceIngestor(address string) (port.PriceIngestor, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := proto.NewPriceIngestorClient(conn)

	return &GRPCPriceIngestor{
		client: client,
		conn:   conn,
	}, nil
}

// IngestPrice implements the PriceIngestor interface
func (g *GRPCPriceIngestor) IngestPrice(ctx context.Context, price *shared_domain.Price) error {
	// Convert domain.Price to proto.PriceMessage
	priceMsg := &proto.PriceMessage{
		Symbol:    price.Symbol,
		Price:     price.Price,
		Currency:  price.Currency,
		Timestamp: price.Timestamp.Unix(),
	}

	// Make the gRPC call
	response, err := g.client.IngestPrice(ctx, priceMsg)
	if err != nil {
		return fmt.Errorf("failed to ingest price via gRPC: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("gRPC ingest failed: %s", response.Message)
	}

	return nil
}

// Close closes the gRPC connection
func (g *GRPCPriceIngestor) Close() error {
	if g.conn != nil {
		return g.conn.Close()
	}
	return nil
}
