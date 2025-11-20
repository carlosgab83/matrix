package domain

import shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"

type Config struct {
	shared_domain.CommonConfig
	DatabaseConnectionString string `json:"database_connection_string" env:"MATRIX_TANK_DATABASE_CONNECTION_STRING"`
	KafkaProducerAddress     string `json:"kafka_producer_address" env:"MATRIX_TANK_KAFKA_PRODUCER_ADDRESS"`
}
