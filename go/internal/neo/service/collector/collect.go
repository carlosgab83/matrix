package collector

import (
	"context"
	"sync"
	"time"

	"github.com/carlosgab83/matrix/go/internal/neo/domain"
	"github.com/carlosgab83/matrix/go/internal/neo/integration/ingestion"
	"github.com/carlosgab83/matrix/go/internal/neo/integration/symbol_fetch"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

type Collector struct {
	Config        domain.Config
	Logger        logging.Logger
	Ctx           context.Context
	Cancel        context.CancelFunc
	Ingestor      ingestion.Ingestor
	SymbolFetcher symbol_fetch.SymbolFetcher
}

func NewCollector(cfg domain.Config, logger logging.Logger, ingestor ingestion.Ingestor, symbolFetcher symbol_fetch.SymbolFetcher) *Collector {
	ctx, cancel := context.WithCancel(context.Background())

	for _, sym := range cfg.Symbols {
		if sym.FetchIntervalSeconds == 0 {
			sym.FetchIntervalSeconds = cfg.DefaultFetchIntervalSeconds
		}
	}

	return &Collector{
		Config:        cfg,
		Logger:        logger,
		Ctx:           ctx,
		Cancel:        cancel,
		Ingestor:      ingestor,
		SymbolFetcher: symbolFetcher,
	}
}

func (c *Collector) Collect() {
	var wg sync.WaitGroup
	buffer := make(chan domain.Symbol, c.Config.WorkersCount*2)

	wg.Add(len(c.Config.Symbols))
	for _, sym := range c.Config.Symbols {
		go c.startSymbolTicker(sym, buffer, &wg)
	}

	for i := 0; i < c.Config.WorkersCount; i++ {
		go c.startSymbolWorker(buffer)
	}

	<-c.Ctx.Done()
	c.Logger.Info("Collector stopped")
	close(buffer)
}

func (c *Collector) startSymbolTicker(sym domain.Symbol, buffer chan<- domain.Symbol, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Duration(sym.FetchIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case buffer <- sym:
				c.Logger.Debug("Tick: Fetch price for symbol", "symbol", sym.Nemo)
			case <-c.Ctx.Done():
				c.Logger.Info("Stopping ticker for symbol", "symbol", sym.Nemo)
				return
			}
		case <-c.Ctx.Done():
			c.Logger.Info("Stopping ticker for symbol", "symbol", sym.Nemo)
			return
		}
	}
}

func (c *Collector) startSymbolWorker(buffer <-chan domain.Symbol) {
	for sym := range buffer {
		c.Logger.Debug("SymbolWorker - Received Tick", "symbol", sym.Nemo)
		var price *shared_domain.Price
		var err error

		switch sym.Nemo {
		case "BTCUSD":
			price, err = c.SymbolFetcher.BTCUSDFetch(c.Ctx)
		case "ETHUSD":
			price, err = c.SymbolFetcher.ETHUSDFetch(c.Ctx)
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
		if err := c.Ingestor.IngestPrice(c.Ctx, price); err != nil {
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
