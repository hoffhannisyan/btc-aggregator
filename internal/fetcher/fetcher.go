package fetcher

import (
	"context"
	"time"
)

type PriceResult struct {
	Source    string
	Price     float64
	FetchedAt time.Time
}

type PriceFetcher interface {
	// FetchPrice fetches the current BTC/USD price
	FetchPrice(ctx context.Context) (*PriceResult, error)
	// Name returns the source name (e.g. "coinbase", "kraken")
	Name() string
}
