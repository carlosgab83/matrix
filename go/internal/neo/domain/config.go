package domain

import shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"

type Config struct {
	shared_domain.CommonConfig
	IngestorAddress             string   `json:"ingestor_address" env:"MATRIX_NEO_INGESTOR_ADDRESS"`
	DefaultFetchIntervalSeconds int      `json:"default_fetch_interval_seconds" env:"MATRIX_NEO_DEFAULT_FETCH_INTERVAL_SECONDS"`
	WorkersCount                int      `json:"workers_count" env:"MATRIX_NEO_WORKERS_COUNT"`
	Symbols                     []Symbol `json:"symbols"`
	GRPCSharedToken             string   `json:"grpc_shared_token" env:"MATRIX_NEO_GRPC_SHARED_TOKEN"`
}

type Symbol struct {
	Nemo                 string `json:"nemo"`
	Name                 string `json:"name"`
	FetchIntervalSeconds int    `json:"fetch_interval_seconds"`
}
