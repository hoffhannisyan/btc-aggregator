package poller

import (
	"context"
	"log/slog"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/hoffhannisyan/btc-aggregator/internal/aggregator"
	"github.com/hoffhannisyan/btc-aggregator/internal/config"
	"github.com/hoffhannisyan/btc-aggregator/internal/fetcher"
	"github.com/hoffhannisyan/btc-aggregator/internal/metrics"
)

type Poller struct {
	fetchers   []fetcher.PriceFetcher
	aggregator *aggregator.Aggregator
	metrics    *metrics.Metrics
	config     *config.Config
	logger     *slog.Logger
}

func New(
	cfg *config.Config,
	agg *aggregator.Aggregator,
	m *metrics.Metrics,
	logger *slog.Logger,
) *Poller {
	client := &http.Client{Timeout: cfg.RequestTimeout}

	fetchers := []fetcher.PriceFetcher{
		fetcher.NewCoinbaseFetcher(client),
		fetcher.NewKrakenFetcher(client),
		fetcher.NewCoinDeskFetcher(client),
	}

	return &Poller{
		fetchers:   fetchers,
		aggregator: agg,
		metrics:    m,
		config:     cfg,
		logger:     logger,
	}
}

func (p *Poller) Start(ctx context.Context) {
	p.logger.Info("poller started", "interval", p.config.PollInterval)

	p.poll(ctx)

	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("poller stopped")
			return
		case <-ticker.C:
			p.poll(ctx)
		}
	}
}

func (p *Poller) poll(ctx context.Context) {
	pollCtx, cancel := context.WithTimeout(ctx, p.config.PollInterval)
	defer cancel()

	var (
		mu     sync.Mutex
		prices []float64
		wg     sync.WaitGroup
	)

	for _, f := range p.fetchers {
		wg.Add(1)
		go func(f fetcher.PriceFetcher) {
			defer wg.Done()

			start := time.Now()
			result, err := p.fetchWithRetry(pollCtx, f)
			latency := time.Since(start)

			if err != nil {
				p.metrics.FetchFailure.WithLabelValues(f.Name()).Inc()
				p.metrics.SourceStatus.WithLabelValues(f.Name()).Set(0)
				p.aggregator.SetSourceHealth(f.Name(), false)
				p.logger.Error("fetch failed",
					"source", f.Name(),
					"error", err,
					"error_type", fetcher.ClassifyError(err),
					"latency", latency,
				)
				return
			}

			p.metrics.FetchSuccess.WithLabelValues(f.Name()).Inc()
			p.metrics.SourceStatus.WithLabelValues(f.Name()).Set(1)
			p.aggregator.SetSourceHealth(f.Name(), true)

			mu.Lock()
			prices = append(prices, result.Price)
			mu.Unlock()

			p.logger.Info("fetch success",
				"source", result.Source,
				"price", result.Price,
				"latency", latency,
			)
		}(f)
	}

	wg.Wait()

	agg := p.aggregator.Aggregate(prices)
	if agg != nil {
		p.metrics.CurrentPrice.Set(agg.Price)
		p.logger.Info("price aggregated",
			"price", agg.Price,
			"sources_used", agg.SourcesUsed,
			"stale", agg.Stale,
		)
	}
}

func (p *Poller) fetchWithRetry(ctx context.Context, f fetcher.PriceFetcher) (*fetcher.PriceResult, error) {
	var lastErr error

	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second

			p.logger.Warn("retrying fetch",
				"source", f.Name(),
				"attempt", attempt,
				"backoff", backoff,
			)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		result, err := f.FetchPrice(ctx)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}

	return nil, lastErr
}
