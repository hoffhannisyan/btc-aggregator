package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const coindeskURL = "https://data-api.coindesk.com/index/cc/v1/latest/tick?market=ccix&instruments=BTC-USD"

type coindeskResponse struct {
	Data map[string]struct {
		Value float64 `json:"VALUE"`
	} `json:"Data"`
}

type CoinDeskFetcher struct {
	client *http.Client
	url    string
}

func NewCoinDeskFetcher(client *http.Client) *CoinDeskFetcher {
	return &CoinDeskFetcher{client: client, url: coindeskURL}
}

func (f *CoinDeskFetcher) Name() string {
	return "coindesk"
}

func (f *CoinDeskFetcher) FetchPrice(ctx context.Context) (*PriceResult, error) {
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

	var result coindeskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	btc, ok := result.Data["BTC-USD"]
	if !ok || btc.Value == 0 {
		return nil, fmt.Errorf("missing BTC-USD data in response")
	}

	return &PriceResult{
		Source:    f.Name(),
		Price:     btc.Value,
		FetchedAt: time.Now(),
	}, nil
}
