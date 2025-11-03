package symbol

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
)

// KrakenResponse represents the structure of Kraken API response
type KrakenResponse struct {
	Result map[string]KrakenPairData `json:"result"`
}

// KrakenPairData represents the data for a trading pair in Kraken
type KrakenPairData struct {
	A []string `json:"a"` // Ask array: [price, whole_lot_volume, lot_volume]
}

// FetchETHUSDPrice fetches the current ETH/USD price from Kraken
func FetchETHUSDPrice(ctx context.Context) (*shared_domain.Price, error) {
	// Create request with context for timeout/cancellation
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.kraken.com/0/public/Ticker?pair=ETHUSD", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse JSON response
	var krakenResp KrakenResponse
	if err := json.NewDecoder(resp.Body).Decode(&krakenResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Navigate to result.XETHZUSD.a.0
	pairData, exists := krakenResp.Result["XETHZUSD"]
	if !exists {
		return nil, fmt.Errorf("XETHZUSD pair not found in response")
	}

	if len(pairData.A) == 0 {
		return nil, fmt.Errorf("'a' array is empty in response")
	}

	// Convert first element of 'a' array (ask price) to float
	price, err := strconv.ParseFloat(pairData.A[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price '%s': %w", pairData.A[0], err)
	}

	// Create Price object
	priceObj := &shared_domain.Price{
		Symbol:    "ETHUSD",
		Price:     price,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	return priceObj, nil
}
