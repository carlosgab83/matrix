package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlosgab83/matrix/go/internal/morpheus/domain"
	"github.com/carlosgab83/matrix/go/internal/morpheus/integration/ingestion"
	"github.com/carlosgab83/matrix/go/internal/morpheus/integration/persisence"
	"github.com/carlosgab83/matrix/go/internal/morpheus/integration/publication"
	"github.com/carlosgab83/matrix/go/internal/morpheus/service"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/configuration"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	matrix_proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type App struct {
	Name            string
	Config          domain.Config
	Logger          logging.Logger
	PriceRepository persisence.PriceRepository
}

const (
	AppName    = "Morpheus"
	ConfigFile = "config/morpheus.json"
)

func NewApp() (*App, error) {
	var cfg domain.Config
	if err := configuration.LoadConfig(&cfg, AppName, ConfigFile); err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	logger, err := logging.NewLogger(cfg.CommonConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %v", err)
	}

	priceRepository, err := persisence.NewPriceRepository(cfg.DatabaseConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error creating repository: %v", err)
	}

	logger.Info("Application started", "app", AppName)

	return &App{
		Name:            AppName,
		Config:          cfg,
		Logger:          logger,
		PriceRepository: priceRepository,
	}, nil
}

func (app *App) Run() {
	app.Logger.Debug("Running...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create Publisher
	publisher, err := publication.NewPublisher(app.Config, app.Logger)
	if err != nil {
		app.Logger.Error("Failed to create publisher", "error", err)
		cancel()
		return
	}

	// Create the domain service (inject dependencies via constructor)
	ingestorService := service.NewIngestorService(ctx, app.Logger, app.PriceRepository, publisher)

	// Create the gRPC server adapter (inject domain service)
	grpcServerAdapter := ingestion.NewGRPCPriceIngestorServer(ctx, ingestorService, app.Logger)

	// Set up gRPC infrastructure
	listener, err := net.Listen("tcp", app.Config.IngestorAddress)
	if err != nil {
		app.Logger.Error("Failed to listen on ingestor address", "address", app.Config.IngestorAddress, "error", err)
		cancel()
		return
	}

	// Create gRPC server with authentication interceptor
	authInterceptor := ingestion.AuthStreamInterceptor(app.Config.GRPCSharedToken, app.Logger)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(authInterceptor))
	matrix_proto.RegisterPriceIngestorServer(grpcServer, grpcServerAdapter)

	log.Println("gRPC server listening on :50051")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Error starting gRPC server: %v", err)
		}
	}()

	sig := <-sigChan
	app.Logger.Info("Received signal, shutting down gracefully", "signal", sig)
	grpcServer.GracefulStop()
	app.Logger.Info("Application stopped")
	app.Logger.Close()
	app.PriceRepository.Close()
}
