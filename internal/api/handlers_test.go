package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hoffhannisyan/btc-aggregator/internal/aggregator"
)

func TestHandlePriceNoData(t *testing.T) {
	agg := aggregator.New()
	handler := NewHandler(agg)

	req := httptest.NewRequest(http.MethodGet, "/price", nil)
	w := httptest.NewRecorder()

	handler.handlePrice(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestHandlePriceWithData(t *testing.T) {
	agg := aggregator.New()
	agg.Aggregate([]float64{67000.0, 68000.0})

	handler := NewHandler(agg)

	req := httptest.NewRequest(http.MethodGet, "/price", nil)
	w := httptest.NewRecorder()

	handler.handlePrice(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if resp["price"] != 67500.0 {
		t.Errorf("expected 67500.0, got %v", resp["price"])
	}
	if resp["currency"] != "USD" {
		t.Errorf("expected USD, got %v", resp["currency"])
	}
	if resp["stale"] != false {
		t.Errorf("expected stale=false, got %v", resp["stale"])
	}
}

func TestHandleHealthNoSources(t *testing.T) {
	agg := aggregator.New()
	handler := NewHandler(agg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.handleHealth(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestHandleHealthWithHealthySource(t *testing.T) {
	agg := aggregator.New()
	agg.SetSourceHealth("coinbase", true)

	handler := NewHandler(agg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandleHealthAllFailing(t *testing.T) {
	agg := aggregator.New()
	agg.SetSourceHealth("coinbase", false)
	agg.SetSourceHealth("kraken", false)

	handler := NewHandler(agg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.handleHealth(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}