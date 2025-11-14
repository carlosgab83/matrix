package ingestion

import (
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
)

// AuthStreamInterceptor validates the shared token for streaming RPC
func AuthStreamInterceptor(sharedToken string, logger logging.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Extract metadata from context
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			logger.Warn("missing metadata in gRPC request")
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// Get authorization header
		tokens := md.Get(authorizationHeader)
		if len(tokens) == 0 {
			logger.Warn("missing authorization token in gRPC request")
			return status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		// Validate token
		clientToken := tokens[0]
		if clientToken != sharedToken {
			logger.Warn("invalid authorization token", "received", clientToken)
			return status.Errorf(codes.Unauthenticated, "invalid authorization token")
		}

		// Token is valid, proceed with the request
		logger.Debug("gRPC request authenticated successfully")
		return handler(srv, ss)
	}
}
