package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	FetchSuccess *prometheus.CounterVec
	FetchFailure *prometheus.CounterVec
	CurrentPrice prometheus.Gauge
	SourceStatus *prometheus.GaugeVec
}

func New() *Metrics {
	m := &Metrics{
		FetchSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fetch_success_total",
				Help: "Total number of successful price fetches per source",
			},
			[]string{"source"},
		),
		FetchFailure: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "fetch_failure_total",
				Help: "Total number of failed price fetches per source",
			},
			[]string{"source"},
		),
		CurrentPrice: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "current_price",
				Help: "Current aggregated BTC/USD price",
			},
		),
		SourceStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "source_status",
				Help: "Current status of each source (1=healthy, 0=unhealthy)",
			},
			[]string{"source"},
		),
	}

	prometheus.MustRegister(m.FetchSuccess)
	prometheus.MustRegister(m.FetchFailure)
	prometheus.MustRegister(m.CurrentPrice)
	prometheus.MustRegister(m.SourceStatus)

	return m
}
