package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const krakenURL = "https://api.kraken.com/0/public/Ticker?pair=XBTUSD"

type krakenResponse struct {
	Error  []string                      `json:"error"`
	Result map[string]krakenTickerResult `json:"result"`
}

type krakenTickerResult struct {
	Close []string `json:"c"`
}

type KrakenFetcher struct {
	client *http.Client
	url    string
}

func NewKrakenFetcher(client *http.Client) *KrakenFetcher {
	return &KrakenFetcher{client: client, url: krakenURL}
}

func (f *KrakenFetcher) Name() string {
	return "kraken"
}

func (f *KrakenFetcher) FetchPrice(ctx context.Context) (*PriceResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result krakenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if len(result.Error) > 0 {
		return nil, fmt.Errorf("kraken API error: %s", result.Error[0])
	}

	ticker, ok := result.Result["XXBTZUSD"]
	if !ok || len(ticker.Close) == 0 {
		return nil, fmt.Errorf("missing XXBTZUSD data in response")
	}

	price, err := strconv.ParseFloat(ticker.Close[0], 64)
	if err != nil {
		return nil, fmt.Errorf("parsing price: %w", err)
	}

	return &PriceResult{
		Source:    f.Name(),
		Price:     price,
		FetchedAt: time.Now(),
	}, nil
}