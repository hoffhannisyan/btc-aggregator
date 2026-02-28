package aggregator

import (
	"testing"
)

func TestMedianOddCount(t *testing.T) {
	prices := []float64{100.0, 200.0, 300.0}
	result := calcMedian(prices)

	if result != 200.0 {
		t.Errorf("expected 200.0, got %f", result)
	}
}

func TestMedianEvenCount(t *testing.T) {
	prices := []float64{100.0, 200.0, 300.0, 400.0}
	result := calcMedian(prices)

	if result != 250.0 {
		t.Errorf("expected 250.0, got %f", result)
	}
}

func TestMedianSinglePrice(t *testing.T) {
	prices := []float64{42000.0}
	result := calcMedian(prices)

	if result != 42000.0 {
		t.Errorf("expected 42000.0, got %f", result)
	}
}

func TestMedianUnsorted(t *testing.T) {
	prices := []float64{300.0, 100.0, 200.0}
	result := calcMedian(prices)

	if result != 200.0 {
		t.Errorf("expected 200.0, got %f", result)
	}
}

func TestAggregateWithPrices(t *testing.T) {
	agg := New()
	result := agg.Aggregate([]float64{67000.0, 68000.0})

	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Price != 67500.0 {
		t.Errorf("expected 67500.0, got %f", result.Price)
	}
	if result.SourcesUsed != 2 {
		t.Errorf("expected 2 sources, got %d", result.SourcesUsed)
	}
	if result.Stale {
		t.Error("expected stale=false")
	}
}

func TestAggregateEmpty_NoLastKnown(t *testing.T) {
	agg := New()
	result := agg.Aggregate([]float64{})

	if result != nil {
		t.Error("expected nil when no prices and no last known")
	}
}

func TestAggregateEmpty_WithLastKnown(t *testing.T) {
	agg := New()

	// first call sets last known
	agg.Aggregate([]float64{67000.0})

	// second call with empty prices should return stale
	result := agg.Aggregate([]float64{})

	if result == nil {
		t.Fatal("expected stale result, got nil")
	}
	if result.Price != 67000.0 {
		t.Errorf("expected 67000.0, got %f", result.Price)
	}
	if !result.Stale {
		t.Error("expected stale=true")
	}
	if result.SourcesUsed != 0 {
		t.Errorf("expected 0 sources, got %d", result.SourcesUsed)
	}
}

func TestSourceHealth(t *testing.T) {
	agg := New()

	if agg.HasHealthySource() {
		t.Error("expected no healthy sources initially")
	}

	agg.SetSourceHealth("coinbase", true)
	if !agg.HasHealthySource() {
		t.Error("expected healthy source after setting coinbase")
	}

	agg.SetSourceHealth("coinbase", false)
	if agg.HasHealthySource() {
		t.Error("expected no healthy sources after setting coinbase unhealthy")
	}
}