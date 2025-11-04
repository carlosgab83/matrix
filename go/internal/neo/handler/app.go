package handler

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlosgab83/matrix/go/internal/neo/domain"
	"github.com/carlosgab83/matrix/go/internal/neo/integration/ingestion"
	"github.com/carlosgab83/matrix/go/internal/neo/service/collector"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/configuration"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

type App struct {
	Name   string
	Config domain.Config
	Logger logging.Logger
}

const (
	AppName    = "Neo"
	ConfigFile = "config/neo.json"
)

func NewApp() (*App, error) {
	var cfg domain.Config
	if err := configuration.LoadConfig(&cfg, AppName, ConfigFile); err != nil {
		log.Fatalf("Error loading config: %v\n", err)
		return nil, err
	}

	if cfg.WorkersCount <= 0 {
		cfg.WorkersCount = 100
	}

	logger, err := logging.NewLogger(cfg.CommonConfig)
	if err != nil {
		log.Fatalf("Error creating logger: %v\n", err)
		return nil, err
	}

	logger.Info("Application started", "app", AppName)

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

	ingestor, err := ingestion.NewIngestor(app.Config)
	if err != nil {
		app.Logger.Error("Failed to create Ingestor client", "error", err)
		return
	}
	defer ingestor.Close()

	coll := collector.NewCollector(app.Config, app.Logger, ingestor)
	go coll.Collect()

	sig := <-sigChan
	app.Logger.Info("Received signal, shutting down gracefully", "signal", sig)
	coll.Stop()
	app.Logger.Info("Application stopped")
	app.Logger.Close()
}
