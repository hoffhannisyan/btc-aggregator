package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const coinbaseURL = "https://api.coinbase.com/v2/prices/BTC-USD/spot"

type coinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

type CoinbaseFetcher struct {
	client *http.Client
	url    string
}

func NewCoinbaseFetcher(client *http.Client) *CoinbaseFetcher {
	return &CoinbaseFetcher{client: client, url: coinbaseURL}
}

func (f *CoinbaseFetcher) Name() string {
	return "coinbase"
}

func (f *CoinbaseFetcher) FetchPrice(ctx context.Context) (*PriceResult, error) {
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

	var result coinbaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	price, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing price: %w", err)
	}

	return &PriceResult{
		Source:    f.Name(),
		Price:     price,
		FetchedAt: time.Now(),
	}, nil
}