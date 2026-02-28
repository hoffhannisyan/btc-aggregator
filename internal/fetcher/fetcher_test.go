package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCoinbaseFetcherSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"amount":"67500.50"}}`))
	}))
	defer server.Close()

	f := &CoinbaseFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	result, err := f.FetchPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Price != 67500.50 {
		t.Errorf("expected 67500.50, got %f", result.Price)
	}
	if result.Source != "coinbase" {
		t.Errorf("expected source coinbase, got %s", result.Source)
	}
}

func TestCoinbaseFetcherBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	f := &CoinbaseFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	_, err := f.FetchPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCoinbaseFetcherBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	}))
	defer server.Close()

	f := &CoinbaseFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	_, err := f.FetchPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestKrakenFetcherSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":[],"result":{"XXBTZUSD":{"c":["67800.00","1.000"]}}}`))
	}))
	defer server.Close()

	f := &KrakenFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	result, err := f.FetchPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Price != 67800.00 {
		t.Errorf("expected 67800.00, got %f", result.Price)
	}
}

func TestKrakenFetcherAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"error":["EGeneral:Internal error"],"result":{}}`))
	}))
	defer server.Close()

	f := &KrakenFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	_, err := f.FetchPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCoinDeskFetcherSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Data":{"BTC-USD":{"VALUE":67900.50}}}`))
	}))
	defer server.Close()

	f := &CoinDeskFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	result, err := f.FetchPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Price != 67900.50 {
		t.Errorf("expected 67900.50, got %f", result.Price)
	}
	if result.Source != "coindesk" {
		t.Errorf("expected source coindesk, got %s", result.Source)
	}
}

func TestCoinDeskFetcherMissingData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Data":{}}`))
	}))
	defer server.Close()

	f := &CoinDeskFetcher{
		client: server.Client(),
		url:    server.URL,
	}

	_, err := f.FetchPrice(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

