package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// KrakenTickerResult represents a subset of Kraken ticker API response.
type KrakenTickerResult struct {
	C []string `json:"c"` // last trade closed: [<price>, <lot volume>]
}

type krakenAPIResponse struct {
	Error  []string                      `json:"error"`
	Result map[string]KrakenTickerResult `json:"result"`
}

type TickerClient struct {
	HTTPClient *http.Client
	BaseURL    string
	CacheTTL   time.Duration
	cache      map[string]cacheEntry
}

type cacheEntry struct {
	price     float64
	fetchedAt time.Time
}

func NewTickerClient() *TickerClient {
	return &TickerClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    "https://api.kraken.com/0/public/Ticker",
		CacheTTL:   time.Minute,
		cache:      make(map[string]cacheEntry),
	}
}

// mapToKrakenPair converts standard pair format to Kraken API format
func mapToKrakenPair(pair string) string {
	switch pair {
	default:
		return strings.ReplaceAll(pair, "/", "")
	}
}

// FetchLastPrices fetches last trade closed price for given pairs from Kraken, with per-minute caching.
func (c *TickerClient) FetchLastPrices(pairs []string) (map[string]float64, error) {
	out := make(map[string]float64)
	var toQuery []string
	var pairMapping = make(map[string]string) // original pair -> kraken query
	
	cutoff := time.Now().Add(-c.CacheTTL)
	for _, p := range pairs {
		if ce, ok := c.cache[p]; ok && ce.fetchedAt.After(cutoff) {
			out[p] = ce.price
			continue
		}
		krakenPair := mapToKrakenPair(p)
		toQuery = append(toQuery, krakenPair)
		pairMapping[p] = krakenPair
	}
	if len(toQuery) == 0 {
		return out, nil
	}

	// Build request
	q := url.Values{}
	q.Set("pair", strings.Join(toQuery, ","))
	endpoint := c.BaseURL + "?" + q.Encode()
	resp, err := c.HTTPClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("kraken request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kraken non-200 status: %d", resp.StatusCode)
	}

	var apiResp krakenAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}
	if len(apiResp.Error) > 0 {
		return nil, fmt.Errorf("kraken errors: %v", apiResp.Error)
	}

	now := time.Now()
	for krakenResponsePair, data := range apiResp.Result {
		if len(data.C) >= 1 {
			priceStr := data.C[0]
			var price float64
			if _, err := fmt.Sscanf(priceStr, "%f", &price); err == nil {
				for originalPair, krakenQuery := range pairMapping {
					// Handle Kraken's various pair naming conventions
					if krakenResponsePair == krakenQuery ||
					   (strings.Contains(krakenResponsePair, "XXBTZ") && strings.Contains(krakenQuery, "BTC")) ||
					   (strings.Contains(krakenResponsePair, "XBT") && strings.Contains(krakenQuery, "BTC")) {
						c.cache[originalPair] = cacheEntry{price: price, fetchedAt: now}
						out[originalPair] = price
						break
					}
				}
			}
		}
	}
	return out, nil
}
