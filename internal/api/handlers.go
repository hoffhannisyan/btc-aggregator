package api

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/hoffhannisyan/btc-aggregator/internal/aggregator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	aggregator *aggregator.Aggregator
}

func NewHandler(agg *aggregator.Aggregator) *Handler {
	return &Handler{aggregator: agg}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/price", h.handlePrice)
	mux.HandleFunc("/health", h.handleHealth)
	mux.Handle("/metrics", promhttp.Handler())
}

func (h *Handler) handlePrice(w http.ResponseWriter, r *http.Request) {
	current := h.aggregator.GetCurrent()
	if current == nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error":"no price data available"}`, http.StatusServiceUnavailable)
		return
	}

	// round to 2 decimal places
	price := math.Round(current.Price*100) / 100

	resp := map[string]interface{}{
		"price":        price,
		"currency":     current.Currency,
		"sources_used": current.SourcesUsed,
		"last_updated": current.LastUpdated.UTC().Format(time.RFC3339),
		"stale":        current.Stale,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if h.aggregator.HasHealthySource() {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
}
