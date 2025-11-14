package domain

import shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"

type Config struct {
	shared_domain.CommonConfig
	IngestorAddress          string `json:"ingestor_address" env:"MATRIX_MORPHEUS_INGESTOR_ADDRESS"`
	DatabaseConnectionString string `json:"database_connection_string" env:"MATRIX_MORPHEUS_DATABASE_CONNECTION_STRING"`
	EventBusID               int    `json:"event_bus_id" env:"MATRIX_MORPHEUS_EVENT_BUS_ID"`
	GRPCSharedToken          string `json:"grpc_shared_token" env:"MATRIX_MORPHEUS_GRPC_SHARED_TOKEN"`
}
