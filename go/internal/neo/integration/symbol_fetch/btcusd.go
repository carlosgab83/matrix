package symbol_fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

type BitstampResponse struct {
	Last string `json:"last"`
}

type BTCUSDFetcher struct {
}

// BTCUSDFetch fetches the current BTC/USD price from Bitstamp
func (sf *BTCUSDFetcher) BTCUSDFetch(ctx context.Context) (*shared_domain.Price, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.bitstamp.net/api/ticker/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var bitstampResp BitstampResponse
	if err := json.NewDecoder(resp.Body).Decode(&bitstampResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	price, err := strconv.ParseFloat(bitstampResp.Last, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price '%s': %w", bitstampResp.Last, err)
	}

	priceObj := &shared_domain.Price{
		Symbol:    "BTCUSD",
		Price:     price,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	return priceObj, nil
}
