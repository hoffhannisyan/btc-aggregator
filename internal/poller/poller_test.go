package poller

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"

	"time"

	"github.com/hoffhannisyan/btc-aggregator/internal/config"
	"github.com/hoffhannisyan/btc-aggregator/internal/fetcher"
)

// mockFetcher implements fetcher.PriceFetcher for testing
type mockFetcher struct {
	name      string
	price     float64
	err       error
	callCount atomic.Int32
}

func (m *mockFetcher) Name() string {
	return m.name
}

func (m *mockFetcher) FetchPrice(ctx context.Context) (*fetcher.PriceResult, error) {
	m.callCount.Add(1)
	if m.err != nil {
		return nil, m.err
	}
	return &fetcher.PriceResult{
		Source:    m.name,
		Price:     m.price,
		FetchedAt: time.Now(),
	}, nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func testConfig() *config.Config {
	return &config.Config{
		PollInterval:   10 * time.Second,
		RequestTimeout: 5 * time.Second,
		MaxRetries:     2,
		HTTPPort:       8080,
	}
}

func TestFetchWithRetrySuccess(t *testing.T) {
	mock := &mockFetcher{name: "test", price: 67000.0}
	cfg := testConfig()
	logger := testLogger()

	p := &Poller{
		config: cfg,
		logger: logger,
	}

	result, err := p.fetchWithRetry(context.Background(), mock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Price != 67000.0 {
		t.Errorf("expected 67000.0, got %f", result.Price)
	}
	if mock.callCount.Load() != 1 {
		t.Errorf("expected 1 call, got %d", mock.callCount.Load())
	}
}

func TestFetchWithRetryAllFail(t *testing.T) {
	mock := &mockFetcher{name: "test", err: fmt.Errorf("connection refused")}
	cfg := testConfig()
	logger := testLogger()

	p := &Poller{
		config: cfg,
		logger: logger,
	}

	_, err := p.fetchWithRetry(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// initial attempt + 2 retries = 3
	if mock.callCount.Load() != 3 {
		t.Errorf("expected 3 calls, got %d", mock.callCount.Load())
	}
}

func TestFetchWithRetryContextCancelled(t *testing.T) {
	mock := &mockFetcher{name: "test", err: fmt.Errorf("connection refused")}
	cfg := testConfig()
	logger := testLogger()

	p := &Poller{
		config: cfg,
		logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := p.fetchWithRetry(ctx, mock)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}