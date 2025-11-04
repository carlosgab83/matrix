package handler

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlosgab83/matrix/go/internal/morpheus/domain"
	"github.com/carlosgab83/matrix/go/internal/morpheus/integration/ingestion"
	"github.com/carlosgab83/matrix/go/internal/morpheus/service"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/configuration"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	matrix_proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
	"google.golang.org/grpc"
	// Driver PostgreSQL
	// _ "github.com/lib/pq"
)

type App struct {
	Name   string
	Config domain.Config
	Logger logging.Logger
}

const (
	AppName    = "Morpheus"
	ConfigFile = "config/morpheus.json"
)

func NewApp() (*App, error) {
	var cfg domain.Config
	if err := configuration.LoadConfig(&cfg, AppName, ConfigFile); err != nil {
		log.Fatalf("Error loading config: %v\n", err)
		return nil, err
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

	// Create the domain service (inject dependencies via constructor)
	ingestorService := service.NewIngestorService(app.Logger)

	// Create the gRPC server adapter (inject domain service)
	grpcServerAdapter := ingestion.NewGRPCPriceIngestorServer(ingestorService)

	// Set up gRPC infrastructure
	listener, err := net.Listen("tcp", app.Config.IngestorAddress)
	if err != nil {
		app.Logger.Error("Failed to listen on ingestor address", "address", app.Config.IngestorAddress, "error", err)
		return
	}

	grpcServer := grpc.NewServer()
	matrix_proto.RegisterPriceIngestorServer(grpcServer, grpcServerAdapter)

	log.Println("gRPC server listening on :50051")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Error starting gRPC server: %v", err)
		}
	}()

	sig := <-sigChan
	app.Logger.Info("Received signal, shutting down gracefully", "signal", sig)
	app.Logger.Info("Application stopped")
	app.Logger.Close()
}
