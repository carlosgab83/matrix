package domain

import shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"

type Config struct {
	shared_domain.CommonConfig
	DefaultFetchIntervalSeconds int      `json:"default_fetch_interval_seconds"`
	WorkersCount                int      `json:"workers_count"`
	Symbols                     []Symbol `json:"symbols"`
}

type Symbol struct {
	Nemo                 string `json:"nemo"`
	Name                 string `json:"name"`
	FetchIntervalSeconds int    `json:"fetch_interval_seconds"`
}
