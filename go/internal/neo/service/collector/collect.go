package collector

import (
	"context"
	"time"

	"github.com/carlosgab83/matrix/go/internal/neo/domain"
	"github.com/carlosgab83/matrix/go/internal/neo/port"
	"github.com/carlosgab83/matrix/go/internal/neo/service/collector/symbol"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	shared_port "github.com/carlosgab83/matrix/go/internal/shared/port"
)

type Collector struct {
	Config        domain.Config
	Logger        shared_port.Logger
	Buffer        chan domain.Symbol
	Ctx           context.Context
	Cancel        context.CancelFunc
	PriceIngestor port.PriceIngestor
}

func NewCollector(cfg domain.Config, logger shared_port.Logger, priceIngestor port.PriceIngestor) *Collector {
	buffer := make(chan domain.Symbol, cfg.WorkersCount*2)

	ctx, cancel := context.WithCancel(context.Background())

	for _, sym := range cfg.Symbols {
		if sym.FetchIntervalSeconds == 0 {
			sym.FetchIntervalSeconds = cfg.DefaultFetchIntervalSeconds
		}
	}

	return &Collector{
		Config:        cfg,
		Logger:        logger,
		Buffer:        buffer,
		Ctx:           ctx,
		Cancel:        cancel,
		PriceIngestor: priceIngestor,
	}
}

func (c *Collector) Collect() {
	for _, sym := range c.Config.Symbols {
		go c.startTickerForSymbol(sym)
	}

	for i := 0; i < c.Config.WorkersCount; i++ {
		go c.processSymbol()
	}

	<-c.Ctx.Done()
	c.Logger.Info("Collector stopped")
	close(c.Buffer)
}

func (c *Collector) startTickerForSymbol(sym domain.Symbol) {
	ticker := time.NewTicker(time.Duration(sym.FetchIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Buffer <- sym
			c.Logger.Debug("Tick: Fetch price for symbol", "symbol", sym.Nemo)
		case <-c.Ctx.Done():
			c.Logger.Info("Stopping ticker for symbol", "symbol", sym.Nemo)
			return
		}
	}
}

func (c *Collector) processSymbol() {
	for sym := range c.Buffer {
		c.Logger.Debug("Received symbol from buffer", "symbol", sym.Nemo)
		var price *shared_domain.Price
		var err error

		switch sym.Nemo {
		case "BTCUSD":
			price, err = symbol.FetchBTCUSDPrice(c.Ctx)
		case "ETHUSD":
			price, err = symbol.FetchETHUSDPrice(c.Ctx)
		default:
			c.Logger.Warn("No handler for symbol", "symbol", sym.Nemo)
			continue
		}

		if err != nil {
			c.Logger.Error("Failed to fetch price",
				"symbol", sym.Nemo,
				"error", err)
			continue
		}

		c.Logger.Info("Fetched price",
			"symbol", price.Symbol,
			"price", price.Price,
			"currency", price.Currency,
			"timestamp", price.Timestamp)

		// Send price to gRPC
		if err := c.PriceIngestor.IngestPrice(c.Ctx, price); err != nil {
			c.Logger.Error("Failed to ingest price via gRPC",
				"symbol", price.Symbol,
				"error", err)
			continue
		}

		c.Logger.Debug("Price sent to gRPC successfully",
			"symbol", price.Symbol)
	}
}

func (c *Collector) Stop() {
	c.Logger.Info("Stopping collector...")
	c.Cancel()
}
