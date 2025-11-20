package handler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlosgab83/matrix/go/internal/shared/integration/configuration"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/tank/domain"
)

type App struct {
	Name   string
	Config domain.Config
	Logger logging.Logger
}

const (
	AppName    = "Tank"
	ConfigFile = "config/tank.json"
)

func NewApp() (*App, error) {
	var cfg domain.Config
	if err := configuration.LoadConfig(&cfg, AppName, ConfigFile); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	logger, err := logging.NewLogger(cfg.CommonConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %w", err)
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Logger.Info("Listening Kafka...")

	<-ctx.Done()
	// TODO: Disconnect Kafka, repository and other services
	app.Logger.Info("Context cancelled, shutting down gracefully", "context", ctx.Err())
	app.Logger.Info("Application stopped")
	app.Logger.Close()
}
