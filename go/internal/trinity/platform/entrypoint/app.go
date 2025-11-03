package entrypoint

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	shared_platform "github.com/carlosgab83/matrix/go/internal/shared/platform"
	shared_port "github.com/carlosgab83/matrix/go/internal/shared/port"
	matrix_proto "github.com/carlosgab83/matrix/go/internal/shared/proto/matrix.proto"
	"github.com/carlosgab83/matrix/go/internal/trinity/adapter"
	"github.com/carlosgab83/matrix/go/internal/trinity/domain"
	"github.com/carlosgab83/matrix/go/internal/trinity/service"
	"google.golang.org/grpc"
)

type App struct {
	Name   string
	Config domain.Config
	Logger shared_port.Logger
}

const (
	AppName    = "Trinity"
	ConfigFile = "config/trinity.json"
)

func NewApp() (*App, error) {
	var cfg domain.Config
	if err := shared_platform.LoadConfig(&cfg, ConfigFile); err != nil {
		log.Fatalf("Error loading config: %v\n", err)
		return nil, err
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

	// 1. Create the domain service (core)
	priceService := service.NewPriceIngestorService(app.Logger)

	// 2. Create the gRPC adapter (inject dependency)
	grpcServerAdapter := adapter.NewGRPCPriceIngestorServer(priceService)

	// 3. Set up gRPC infrastructure
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		app.Logger.Error("Failed to listen on port 50051", "error", err)
		return
	}

	grpcServer := grpc.NewServer()
	matrix_proto.RegisterPriceIngestorServer(grpcServer, grpcServerAdapter)

	log.Println("Servidor gRPC escuchando en :50051")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Error al iniciar gRPC server: %v", err)
		}
	}()

	sig := <-sigChan
	app.Logger.Info("Received signal, shutting down gracefully", "signal", sig)
	app.Logger.Info("Application stopped")
	app.Logger.Close()
}
