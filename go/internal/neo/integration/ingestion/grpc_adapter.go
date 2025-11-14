package ingestion

import (
	"context"
	"fmt"
	"io"
	"sync"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// GRPCPriceIngestor implements the PriceIngestor port using gRPC
type GRPCPriceIngestor struct {
	client      proto.PriceIngestorClient
	conn        *grpc.ClientConn
	stream      proto.PriceIngestor_IngestPriceClient
	mutex       sync.Mutex
	sharedToken string
}

// NewGRPCPriceIngestor creates a new gRPC adapter
func NewGRPCPriceIngestor(address string, sharedToken string) (Ingestor, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := proto.NewPriceIngestorClient(conn)

	return &GRPCPriceIngestor{
		client:      client,
		conn:        conn,
		stream:      nil,
		sharedToken: sharedToken,
	}, nil
}

// IngestPrice implements the PriceIngestor interface
func (g *GRPCPriceIngestor) IngestPrice(ctx context.Context, price *shared_domain.Price) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Convert domain.Price to proto.PriceMessage
	priceMsg := &proto.PriceMessage{
		Symbol:    price.Symbol,
		Price:     price.Price,
		Currency:  price.Currency,
		Timestamp: price.Timestamp.Unix(),
	}

	if g.stream == nil {
		// Add authorization token to context metadata
		md := metadata.New(map[string]string{
			"authorization": g.sharedToken,
		})
		ctx = metadata.NewOutgoingContext(ctx, md)

		stream, err := g.client.IngestPrice(ctx)
		if err != nil {
			return fmt.Errorf("failed to create gRPC stream: %w", err)
		}

		g.stream = stream
	}

	// Make the gRPC call
	if err := g.stream.Send(priceMsg); err != nil {
		g.stream = nil
		return fmt.Errorf("gRPC ingest failed: %w", err)
	}

	return nil
}

// Close closes the gRPC connection
func (g *GRPCPriceIngestor) Close() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var streamErr, connErr error

	if g.stream != nil {
		_, err := g.stream.CloseAndRecv()
		if err != nil && err != io.EOF {
			streamErr = fmt.Errorf("failed to close gRPC stream: %w", err)
		}
	}

	if g.conn != nil {
		connErr = g.conn.Close()
	}

	if streamErr != nil {
		return streamErr
	}

	return connErr
}
