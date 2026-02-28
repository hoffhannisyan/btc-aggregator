package aggregator

import (
	"sort"
	"sync"
	"time"
)

type AggregatedPrice struct {
	Price       float64
	Currency    string
	SourcesUsed int
	LastUpdated time.Time
	Stale       bool
}

type Aggregator struct {
	mu            sync.RWMutex
	lastKnown     *AggregatedPrice
	sourceHealthy map[string]bool
}

func New() *Aggregator {
	return &Aggregator{
		sourceHealthy: make(map[string]bool),
	}
}

func (a *Aggregator) Aggregate(prices []float64) *AggregatedPrice {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(prices) == 0 {
		if a.lastKnown != nil {
			return &AggregatedPrice{
				Price:       a.lastKnown.Price,
				Currency:    "USD",
				SourcesUsed: 0,
				LastUpdated: a.lastKnown.LastUpdated,
				Stale:       true,
			}
		}
		return nil
	}

	median := calcMedian(prices)

	result := &AggregatedPrice{
		Price:       median,
		Currency:    "USD",
		SourcesUsed: len(prices),
		LastUpdated: time.Now(),
		Stale:       false,
	}

	a.lastKnown = result
	return result
}

func (a *Aggregator) GetCurrent() *AggregatedPrice {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.lastKnown
}

func (a *Aggregator) SetSourceHealth(source string, healthy bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sourceHealthy[source] = healthy
}

func (a *Aggregator) HasHealthySource() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, healthy := range a.sourceHealthy {
		if healthy {
			return true
		}
	}
	return false
}

func calcMedian(prices []float64) float64 {
	sorted := make([]float64, len(prices))
	copy(sorted, prices)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}
