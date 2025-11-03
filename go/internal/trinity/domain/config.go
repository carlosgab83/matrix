package domain

import shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"

type Config struct {
	shared_domain.CommonConfig
	DatabaseConnectionString string `json:"database_connection_string"`
	EventBusID               int    `json:"event_bus_id"`
}
