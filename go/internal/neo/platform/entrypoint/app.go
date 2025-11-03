package entrypoint

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlosgab83/matrix/go/internal/neo/adapter"
	"github.com/carlosgab83/matrix/go/internal/neo/domain"
	"github.com/carlosgab83/matrix/go/internal/neo/service/collector"
	shared_platform "github.com/carlosgab83/matrix/go/internal/shared/platform"
	shared_port "github.com/carlosgab83/matrix/go/internal/shared/port"
)

type App struct {
	Name   string
	Config domain.Config
	Logger shared_port.Logger
}

const (
	AppName    = "Neo"
	ConfigFile = "config/neo.json"
)

func NewApp() (*App, error) {
	var cfg domain.Config
	if err := shared_platform.LoadConfig(&cfg, ConfigFile); err != nil {
		log.Fatalf("Error loading config: %v\n", err)
		return nil, err
	}

	if cfg.WorkersCount <= 0 {
		cfg.WorkersCount = 100
	}

	logger, err := shared_platform.NewLogger(cfg.LogFilePath, cfg.LogLevel)
	if err != nil {
		log.Fatalf("Error creating logger: %v\n", err)
		return nil, err
	}

	logger.Logger.Info("Application started", "app", AppName)

	return &App{
		Name:   AppName,
		Config: cfg,
		Logger: logger,
	}, nil
}

func (app *App) Run() {
	app.Logger.Debug("Running...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create gRPC price ingestor
	priceIngestor, err := adapter.NewGRPCPriceIngestor("localhost:50051")
	if err != nil {
		app.Logger.Error("Failed to create gRPC client", "error", err)
		return
	}
	defer priceIngestor.(*adapter.GRPCPriceIngestor).Close()

	// Create collector with dependencies
	coll := collector.NewCollector(app.Config, app.Logger, priceIngestor)
	go coll.Collect()

	sig := <-sigChan
	app.Logger.Info("Received signal, shutting down gracefully", "signal", sig)
	coll.Stop()
	app.Logger.Info("Application stopped")
	app.Logger.Close()
}
